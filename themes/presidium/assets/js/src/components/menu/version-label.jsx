import React from 'react';

const VersionLabel = () => {
    return (
        <li className='presidium-js-version' style={{
            textAlign: 'center',
            color: '#a7a5a5',
            marginTop: '10px',
            fontSize: '12px'
        }}>
            v{PRESIDIUM_VERSION}
        </li>
    )
};

export default VersionLabel;


