import React, { Component } from 'react';
import { getVersions } from '../../actions/versions';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import {path} from '../../util/articles';

/**
 * Version Navigation Component.
 */
class Versions extends Component {

    constructor(props) {
        super(props);

        const siteroot = window.presidium.versions.siteroot;
        const v = window.location.pathname.replace(siteroot, '').replace(/^\/([^\/]*).*$/, '$1');

        this.state = {
            versioned: false,
            versions: {},
            siteroot: siteroot,
            selected_version: v || 'latest'
        };
    }

    onChangeVersion(e) {
        const version =  e.target.value === 'latest' ?  '' : e.target.value;
        window.location.href = path.concat(this.state.siteroot, version);
    }

    componentWillReceiveProps(props) {
        this.setState(props.versions);
    }

    componentWillMount(){
        this.props.getVersions(path.concat(this.state.siteroot, 'versions.json'));
    }

    render() {
        return this.state.versioned && (
            <div className="filter versions-filter form-group">
                <select id="versions-select" className="form-control" value={this.state.selected_version}
                        onChange={(version) => this.onChangeVersion(version)}>
                    {
                        this.state.versions.map(version => {
                            return <option key={ version } value={ version }>v. { version }</option>
                        })
                    }
                </select>
            </div>)
    }
}

function mapStateToProps(state) {
    return { enabled: state.enabled, versions: state.versions };
}

function mapDispatchToProps(dispatch){
    return bindActionCreators({ getVersions }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(Versions);
