import { combineReducers } from 'redux';
import VersionsReducer from './reducer_versions';

const rootReducer = combineReducers({
    versions: VersionsReducer
});

export default rootReducer;