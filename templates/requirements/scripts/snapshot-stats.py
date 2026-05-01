#!/usr/bin/env python3
"""
Discovery Stats Snapshot Generator

Generates time-series statistics snapshots from domain modeling artifacts.
Supports multiple artifact types: YAML files with status tracking, DBML entity models, etc.

Usage:
    python scripts/snapshot-stats.py                    # Add new snapshot
    python scripts/snapshot-stats.py --dry-run          # Preview without saving
    python scripts/snapshot-stats.py --config custom.yaml  # Use custom config
"""

import argparse
import re
import sys
from datetime import datetime, timezone
from pathlib import Path

import yaml


# Default configuration - can be overridden with --config
DEFAULT_CONFIG = {
    "output_file": "data/discovery-stats.yaml",
    "categories": [
        {"name": "domain_model", "label": "Domain Model", "order": 1},
        {"name": "requirements", "label": "Requirements", "order": 2},
        {"name": "tests", "label": "Tests", "order": 3},
    ],
    "artifacts": [
        # Domain Model artifacts
        {
            "name": "archetypes",
            "source": "agent-artifacts/domain-modeling/user-archetypes.yaml",
            "type": "yaml_list",
            "path": "archetypes",
            "track_status": True,
            "label": "Archetypes",
            "category": "domain_model",
            "order": 1,
        },
        {
            "name": "personas",
            "source": "agent-artifacts/domain-modeling/user-archetypes.yaml",
            "type": "yaml_nested_list",
            "parent_path": "archetypes",
            "child_key": "personas",
            "track_status": True,
            "label": "Personas",
            "category": "domain_model",
            "order": 2,
        },
        {
            "name": "entities",
            "source": "agent-artifacts/domain-modeling/logical-entity-model.dbml",
            "type": "dbml_tables",
            "track_status": False,
            "label": "Entities",
            "category": "domain_model",
            "order": 3,
        },
        {
            "name": "enums",
            "source": "agent-artifacts/domain-modeling/logical-entity-model.dbml",
            "type": "dbml_enums",
            "track_status": False,
            "label": "Enums",
            "category": "domain_model",
            "order": 4,
        },
        {
            "name": "functional_areas",
            "source": "agent-artifacts/domain-modeling/functional-hierarchy.yaml",
            "type": "yaml_list",
            "path": "functional_areas",
            "track_status": True,
            "label": "Functional Areas",
            "category": "domain_model",
            "order": 5,
        },
        {
            "name": "features",
            "source": "agent-artifacts/domain-modeling/functional-hierarchy.yaml",
            "type": "yaml_nested_list",
            "parent_path": "functional_areas",
            "child_key": "features",
            "track_status": True,
            "label": "Features",
            "category": "domain_model",
            "order": 6,
        },
        {
            "name": "functions",
            "source": "agent-artifacts/domain-modeling/functional-hierarchy.yaml",
            "type": "yaml_deep_nested_list",
            "parent_path": "functional_areas",
            "child_key": "features",
            "grandchild_key": "functions",
            "track_status": True,
            "label": "Functions",
            "category": "domain_model",
            "order": 7,
        },
        {
            "name": "processes",
            "source": "agent-artifacts/domain-modeling/process-flows.yaml",
            "type": "yaml_list",
            "path": "processes",
            "track_status": True,
            "label": "Processes",
            "category": "domain_model",
            "order": 8,
        },
        # Requirements artifacts
        {
            "name": "scenarios",
            "source": "agent-artifacts/requirements-content/scenarios.yaml",
            "type": "yaml_list",
            "path": "scenarios",
            "track_status": True,
            "label": "Scenarios",
            "category": "requirements",
            "order": 1,
        },
        # Test artifacts
        {
            "name": "test_instances",
            "source": "data/test_results.yaml",
            "type": "test_results",
            "path": "spec_files",
            "track_status": True,
            "label": "Functions Tested",
            "category": "tests",
            "order": 1,
        },
        {
            "name": "scenarios_executed",
            "source": "data/test_results.yaml",
            "type": "test_results",
            "path": "test_status",
            "track_status": True,
            "label": "Scenarios Tested",
            "category": "tests",
            "order": 2,
        },
    ],
}


def load_yaml_file(filepath: Path) -> dict:
    """Load and parse a YAML file."""
    with open(filepath, "r", encoding="utf-8") as f:
        return yaml.safe_load(f)


def load_dbml_file(filepath: Path) -> str:
    """Load a DBML file as text."""
    with open(filepath, "r", encoding="utf-8") as f:
        return f.read()


def count_by_status(items: list) -> dict:
    """Count items grouped by their status field."""
    status_counts = {}
    for item in items:
        status = item.get("status", "unknown")
        status_counts[status] = status_counts.get(status, 0) + 1
    return status_counts


def extract_yaml_list(data: dict, path: str) -> list:
    """Extract a list from a YAML structure using a dot-separated path."""
    parts = path.split(".")
    current = data
    for part in parts:
        if current is None:
            return []
        current = current.get(part, [])
    return current if isinstance(current, list) else []


def extract_nested_list(data: dict, parent_path: str, child_key: str) -> list:
    """Extract items from a nested list structure."""
    parents = extract_yaml_list(data, parent_path)
    items = []
    for parent in parents:
        children = parent.get(child_key, [])
        if isinstance(children, list):
            items.extend(children)
    return items


def extract_deep_nested_list(
    data: dict, parent_path: str, child_key: str, grandchild_key: str
) -> list:
    """Extract items from a deeply nested list structure (3 levels)."""
    parents = extract_yaml_list(data, parent_path)
    items = []
    for parent in parents:
        children = parent.get(child_key, [])
        if isinstance(children, list):
            for child in children:
                grandchildren = child.get(grandchild_key, [])
                if isinstance(grandchildren, list):
                    items.extend(grandchildren)
    return items


def count_dbml_tables(content: str) -> int:
    """Count Table definitions in DBML content."""
    # Match "Table name {" or "Table name as alias {"
    pattern = r"^Table\s+\w+"
    matches = re.findall(pattern, content, re.MULTILINE)
    return len(matches)


def count_dbml_enums(content: str) -> int:
    """Count Enum definitions in DBML content."""
    pattern = r"^Enum\s+\w+"
    matches = re.findall(pattern, content, re.MULTILINE)
    return len(matches)


def count_test_results(data: dict, path: str) -> dict:
    """Count test results by status (pass/failed/skipped)."""
    test_status = data.get(path, {})
    status_counts = {}
    for test_id, test_data in test_status.items():
        status = test_data.get("status", "unknown")
        status_counts[status] = status_counts.get(status, 0) + 1
    return {
        "total": len(test_status),
        "by_status": status_counts,
    }


def collect_artifact_stats(artifact_config: dict, base_path: Path) -> dict:
    """Collect statistics for a single artifact based on its configuration."""
    source_path = base_path / artifact_config["source"]
    artifact_type = artifact_config["type"]

    if not source_path.exists():
        return {"total": 0, "error": f"Source file not found: {artifact_config['source']}"}

    stats = {}

    try:
        if artifact_type == "yaml_list":
            data = load_yaml_file(source_path)
            items = extract_yaml_list(data, artifact_config["path"])
            stats["total"] = len(items)
            if artifact_config.get("track_status"):
                stats["by_status"] = count_by_status(items)

        elif artifact_type == "yaml_nested_list":
            data = load_yaml_file(source_path)
            items = extract_nested_list(
                data, artifact_config["parent_path"], artifact_config["child_key"]
            )
            stats["total"] = len(items)
            if artifact_config.get("track_status"):
                stats["by_status"] = count_by_status(items)

        elif artifact_type == "yaml_deep_nested_list":
            data = load_yaml_file(source_path)
            items = extract_deep_nested_list(
                data,
                artifact_config["parent_path"],
                artifact_config["child_key"],
                artifact_config["grandchild_key"],
            )
            stats["total"] = len(items)
            if artifact_config.get("track_status"):
                stats["by_status"] = count_by_status(items)

        elif artifact_type == "dbml_tables":
            content = load_dbml_file(source_path)
            stats["total"] = count_dbml_tables(content)

        elif artifact_type == "dbml_enums":
            content = load_dbml_file(source_path)
            stats["total"] = count_dbml_enums(content)

        elif artifact_type == "test_results":
            data = load_yaml_file(source_path)
            result = count_test_results(data, artifact_config["path"])
            stats["total"] = result["total"]
            if artifact_config.get("track_status"):
                stats["by_status"] = result["by_status"]

        else:
            stats["error"] = f"Unknown artifact type: {artifact_type}"

    except Exception as e:
        stats["error"] = str(e)

    return stats


def generate_snapshot(config: dict, base_path: Path) -> dict:
    """Generate a complete snapshot of all artifact statistics."""
    now = datetime.now(timezone.utc)
    snapshot = {
        "date": now.strftime("%Y-%m-%d"),
        "timestamp": now.isoformat(),
        "categories": {
            cat["name"]: {"label": cat["label"], "order": cat.get("order", 99)}
            for cat in config.get("categories", [])
        },
        "artifacts": {},
    }

    for artifact_config in config["artifacts"]:
        name = artifact_config["name"]
        stats = collect_artifact_stats(artifact_config, base_path)
        stats["label"] = artifact_config.get("label", name)
        stats["category"] = artifact_config.get("category", "other")
        stats["order"] = artifact_config.get("order", 99)
        snapshot["artifacts"][name] = stats

    return snapshot


def load_existing_stats(filepath: Path) -> dict:
    """Load existing stats file or create empty structure."""
    if filepath.exists():
        with open(filepath, "r", encoding="utf-8") as f:
            return yaml.safe_load(f) or {"snapshots": []}
    return {
        "metadata": {
            "title": "Discovery Artifact Statistics",
            "description": "Time-series tracking of domain modeling artifact counts",
            "generated_by": "scripts/snapshot-stats.py",
        },
        "snapshots": [],
    }


def save_stats(data: dict, filepath: Path) -> None:
    """Save stats to YAML file."""
    filepath.parent.mkdir(parents=True, exist_ok=True)
    with open(filepath, "w", encoding="utf-8") as f:
        yaml.dump(data, f, default_flow_style=False, sort_keys=False, allow_unicode=True)


def print_snapshot_summary(snapshot: dict) -> None:
    """Print a human-readable summary of the snapshot."""
    print(f"\n{'='*60}")
    print(f"Snapshot for {snapshot['date']}")
    print(f"{'='*60}")

    # Group artifacts by category
    categories = snapshot.get("categories", {})
    artifacts_by_category = {}
    for name, stats in snapshot["artifacts"].items():
        category = stats.get("category", "other")
        if category not in artifacts_by_category:
            artifacts_by_category[category] = []
        artifacts_by_category[category].append((name, stats))

    # Print by category
    for category_name, category_data in categories.items():
        category_label = category_data["label"] if isinstance(category_data, dict) else category_data
        if category_name in artifacts_by_category:
            print(f"\n  {category_label}:")
            print(f"  {'-' * (len(category_label) + 1)}")
            for name, stats in artifacts_by_category[category_name]:
                label = stats.get("label", name)
                total = stats.get("total", "N/A")

                if "error" in stats:
                    print(f"    {label}: ERROR - {stats['error']}")
                elif "by_status" in stats:
                    status_str = ", ".join(
                        f"{k}: {v}" for k, v in sorted(stats["by_status"].items())
                    )
                    print(f"    {label}: {total} ({status_str})")
                else:
                    print(f"    {label}: {total}")

    # Print any uncategorized artifacts
    if "other" in artifacts_by_category:
        print("\n  Other:")
        print("  ------")
        for name, stats in artifacts_by_category["other"]:
            label = stats.get("label", name)
            total = stats.get("total", "N/A")
            print(f"    {label}: {total}")

    print()


def main():
    parser = argparse.ArgumentParser(
        description="Generate discovery artifact statistics snapshot"
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Preview snapshot without saving",
    )
    parser.add_argument(
        "--config",
        type=Path,
        help="Path to custom configuration YAML file",
    )
    parser.add_argument(
        "--base-path",
        type=Path,
        default=Path.cwd(),
        help="Base path for resolving relative paths (default: current directory)",
    )
    args = parser.parse_args()

    # Load configuration
    if args.config:
        config = load_yaml_file(args.config)
    else:
        config = DEFAULT_CONFIG

    # Generate snapshot
    base_path = args.base_path
    snapshot = generate_snapshot(config, base_path)

    # Print summary
    print_snapshot_summary(snapshot)

    if args.dry_run:
        print("Dry run - no changes saved.")
        return 0

    # Load existing stats and append new snapshot
    output_path = base_path / config["output_file"]
    stats_data = load_existing_stats(output_path)

    # Check if we already have a snapshot for today
    today = snapshot["date"]
    existing_dates = [s.get("date") for s in stats_data.get("snapshots", [])]

    if today in existing_dates:
        print(f"Snapshot for {today} already exists. Updating...")
        stats_data["snapshots"] = [
            s for s in stats_data["snapshots"] if s.get("date") != today
        ]

    stats_data["snapshots"].append(snapshot)

    # Sort snapshots by date
    stats_data["snapshots"].sort(key=lambda s: s.get("date", ""))

    # Save
    save_stats(stats_data, output_path)
    print(f"Snapshot saved to {output_path}")

    return 0


if __name__ == "__main__":
    sys.exit(main())
