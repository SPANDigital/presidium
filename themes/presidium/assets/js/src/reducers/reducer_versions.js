import { GET_VERSIONS } from '../actions/versions';

export default function(state = {}, action) {
    switch(action.type) {
        case GET_VERSIONS:
            return Object.assign({}, state, action.payload.data);
        default:
            return state;
    }
}