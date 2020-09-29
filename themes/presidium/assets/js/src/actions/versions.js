import axios from 'axios';
export const GET_VERSIONS = 'GET_VERSIONS';

export function getVersions(url) {
    const request = axios.get(url);
    return {
        type: GET_VERSIONS,
        payload: request
    };
}
