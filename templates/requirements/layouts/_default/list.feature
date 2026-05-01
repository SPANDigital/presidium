{{- if eq .Params.type "story" -}}
{{- $filename := printf "%s.feature" .File.BaseFileName -}}
{{- $story := .Content | htmlUnescape | replaceRE "<[^>]*>" "" | replaceRE "`([^`]+)`" "$1" -}}
{{- $story = replaceRE "(?m)\n+$" "" $story -}}
{{/* ---------- Add user archetype tags if present ---------- */}}
{{- with .Params.archetypes }}{{- if gt (len .) 0 }}{{- range $index, $archetype := . }}{{ if $index }} {{ end }}@user:{{ replace $archetype " " "_" }}{{ end }}
{{ end }}{{- end -}}
{{- $story | safeHTML }}

{{/* ---------- Find scenarios under current path ---------- */}}
{{- $currentPath := .File.Dir -}}
{{- $scenarios := where .Site.Pages "Params.type" "scenario" -}}
{{- $matchingScenarios := slice -}}
{{- range $scenarios -}}
  {{- $scenarioDir := .File.Dir -}}
  {{/* Only include scenarios that are direct children (same directory) */}}
  {{- if eq $scenarioDir $currentPath -}}
    {{- $matchingScenarios = $matchingScenarios | append . -}}
  {{- end -}}
{{- end -}}

{{/* ---------- Create site-wide scenario map for dependency lookup ---------- */}}
{{- $scenarioMap := dict -}}
{{- range $scenarios -}}
{{- if .Params.scenarioId -}}
{{- $scenarioMap = merge $scenarioMap (dict (string .Params.scenarioId) .) -}}
{{- end -}}
{{- end -}}

{{/* ---------- Build ordered list with transitive dependency resolution ---------- */}}
{{/* Resolves up to 3 levels deep: dep -> dep's dep -> dep's dep's dep */}}
{{/* This covers chains like: create_client -> enroll -> exit -> reopen */}}
{{- $orderedScenarios := slice -}}
{{- $processedIds := slice -}}
{{- $dependencyIds := slice -}}

{{/* Process each scenario and its full transitive dependency chain */}}
{{- range $matchingScenarios -}}
  {{- $scenario := . -}}
  {{- $scenarioId := string .Params.scenarioId -}}

  {{/* Skip if already processed */}}
  {{- if not (in $processedIds $scenarioId) -}}

    {{/* Level 1: direct dependencies */}}
    {{- if .Params.depends_on -}}
      {{- $deps1 := sort .Params.depends_on "sequence_order" -}}
      {{- range $deps1 -}}
        {{- $dep1Id := string (index . "scenario_id") -}}
        {{- $dep1Page := index $scenarioMap $dep1Id -}}

        {{/* Level 2: dependencies of dependencies */}}
        {{- if $dep1Page -}}
          {{- if $dep1Page.Params.depends_on -}}
            {{- $deps2 := sort $dep1Page.Params.depends_on "sequence_order" -}}
            {{- range $deps2 -}}
              {{- $dep2Id := string (index . "scenario_id") -}}
              {{- $dep2Page := index $scenarioMap $dep2Id -}}

              {{/* Level 3: dependencies of dependencies of dependencies */}}
              {{- if $dep2Page -}}
                {{- if $dep2Page.Params.depends_on -}}
                  {{- $deps3 := sort $dep2Page.Params.depends_on "sequence_order" -}}
                  {{- range $deps3 -}}
                    {{- $dep3Id := string (index . "scenario_id") -}}
                    {{- if not (in $processedIds $dep3Id) -}}
                      {{- if index $scenarioMap $dep3Id -}}
                        {{- $orderedScenarios = $orderedScenarios | append (index $scenarioMap $dep3Id) -}}
                        {{- $processedIds = $processedIds | append $dep3Id -}}
                        {{- $dependencyIds = $dependencyIds | append $dep3Id -}}
                      {{- else -}}
                        {{- $orderedScenarios = $orderedScenarios | append (dict "isMissing" true "scenarioId" $dep3Id) -}}
                      {{- end -}}
                    {{- end -}}
                  {{- end -}}
                {{- end -}}
              {{- end -}}

              {{/* Add level 2 dependency */}}
              {{- if not (in $processedIds $dep2Id) -}}
                {{- if $dep2Page -}}
                  {{- $orderedScenarios = $orderedScenarios | append $dep2Page -}}
                  {{- $processedIds = $processedIds | append $dep2Id -}}
                  {{- $dependencyIds = $dependencyIds | append $dep2Id -}}
                {{- else -}}
                  {{- $orderedScenarios = $orderedScenarios | append (dict "isMissing" true "scenarioId" $dep2Id) -}}
                {{- end -}}
              {{- end -}}
            {{- end -}}
          {{- end -}}
        {{- end -}}

        {{/* Add level 1 dependency */}}
        {{- if not (in $processedIds $dep1Id) -}}
          {{- if $dep1Page -}}
            {{- $orderedScenarios = $orderedScenarios | append $dep1Page -}}
            {{- $processedIds = $processedIds | append $dep1Id -}}
            {{- $dependencyIds = $dependencyIds | append $dep1Id -}}
          {{- else -}}
            {{- $orderedScenarios = $orderedScenarios | append (dict "isMissing" true "scenarioId" $dep1Id) -}}
          {{- end -}}
        {{- end -}}
      {{- end -}}
    {{- end -}}

    {{/* Then add the scenario itself */}}
    {{- $orderedScenarios = $orderedScenarios | append $scenario -}}
    {{- $processedIds = $processedIds | append $scenarioId -}}
  {{- end -}}
{{- end -}}

{{- $matchingScenarios = $orderedScenarios -}}

{{/* ---------- Render each scenario ---------- */}}
{{- range $matchingScenarios }}
{{- if not .File -}}
{{/* This is a missing dependency marker (dict), not a page */}}
# Missing dependency: scenario {{ index . "scenarioId" }}

{{- else -}}
{{- $currentScenarioId := string .Params.scenarioId -}}
{{- if in $dependencyIds $currentScenarioId -}}
{{/* This scenario is a dependency - add @dependency tag */}}
@scenario:{{ .Params.scenarioId }}{{- if .Params.testSuites }}{{- range .Params.testSuites }} @test_suite:{{ replace . " " "_" }}{{- end }}{{- end }} @dependency{{ "\n" }}
{{- else -}}
{{/* This is a main scenario - use standard tags */}}
{{- partial "gherkin-tags.html" . }}
{{- end -}}
{{- $sraw := .RawContent -}}
{{- $st := or .Params.title .Title -}}
{{/* Use ### Test section if available, otherwise fall back to full content */}}
{{- $testParts := split $sraw "### Test" -}}
{{- if gt (len $testParts) 1 -}}
{{- $sraw = index $testParts 1 -}}
{{- else -}}
{{/* No ### Test section — strip ### Logical heading if present, use remaining content */}}
{{- $logicalParts := split $sraw "### Logical" -}}
{{- if gt (len $logicalParts) 1 -}}
{{- $sraw = index $logicalParts 1 -}}
{{- end -}}
{{- end -}}
{{- $sraw = replaceRE `\{\{<\s*article-title(?:\s+[^>]+)?\s*>\}\}` (printf "%s" $st) $sraw -}}
{{- $sraw = replaceRE "([^`])`([^`]+)`([^`])" "$1$2$3" $sraw -}}
{{- $sraw = replaceRE "```gherkin\\s*" "" $sraw -}}
{{- $sraw = replaceRE "```\\s*$" "" $sraw -}}
{{- $sraw = replaceRE "\\*\\*([^*]+)\\*\\*" "$1" $sraw -}}
{{- $sraw = replaceRE "^\\s*" "" $sraw -}}
{{- $sraw | htmlUnescape | safeHTML }}
{{- end -}}
{{- end }}
{{- end -}}
