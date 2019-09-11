import {getConfig} from 'mattermost-redux/selectors/entities/general';

import {getUserId, getTeamId} from 'selectors';

import {id as pluginId} from './manifest';

export function postDropdownMenuAction(postId, menuItemId) {
    return async (dispatch, getState) => {
        const state = getState();
        const opts = {
            postId,
            userId: getUserId(state),
            teamId: getTeamId(state),
            timeId: menuItemId.replace('submenu.', ''),
        };

        fetch(getPluginServerRoute(state) + '/remind/' + postId, {
            method: 'post',
            body: JSON.stringify(opts),
        });
    };
}

export const getPluginServerRoute = (state) => {
    const config = getConfig(state);

    let basePath = '/';
    if (config && config.SiteURL) {
        basePath = new URL(config.SiteURL).pathname;

        if (basePath && basePath[basePath.length - 1] === '/') {
            basePath = basePath.substr(0, basePath.length - 1);
        }
    }

    return basePath + '/plugins/' + pluginId;
};
