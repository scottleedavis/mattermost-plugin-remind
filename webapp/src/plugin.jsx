import React from 'react';

import PluginId from './plugin_id';

import Root from './components/root';

// import BottomTeamSidebar from './components/bottom_team_sidebar';
// import LeftSidebarHeader from './components/left_sidebar_header';
// import UserAttributes from './components/user_attributes';
// import UserActions from './components/user_actions';
// import PostType from './components/post_type';
import {

    // MainMenuMobileIcon,
    ChannelHeaderButtonIcon,

    // FileUploadMethodIcon,
} from './components/icons';
import {

    // mainMenuAction,
    channelHeaderButtonAction,

    // fileUploadMethodAction,
    // postDropdownMenuAction,
    // websocketStatusChange,
    getStatus,
} from './actions';
import reducer from './reducer';

export default class DemoPlugin {
    initialize(registry, store) {
        registry.registerRootComponent(Root);

        // registry.registerPopoverUserAttributesComponent(UserAttributes);
        // registry.registerPopoverUserActionsComponent(UserActions);
        // registry.registerLeftSidebarHeaderComponent(LeftSidebarHeader);
        // registry.registerBottomTeamSidebarComponent(
        //     BottomTeamSidebar,
        // );

        registry.registerChannelHeaderButtonAction(
            <ChannelHeaderButtonIcon/>,
            () => store.dispatch(channelHeaderButtonAction()),
            'Remind Plugin',
        );

        //
        // registry.registerPostTypeComponent('custom_demo_plugin', PostType);
        //
        // registry.registerMainMenuAction(
        //     'Demo Plugin',
        //     () => store.dispatch(mainMenuAction()),
        //     <MainMenuMobileIcon/>,
        // );
        //
        // registry.registerPostDropdownMenuAction(
        //     'Demo Plugin',
        //     () => store.dispatch(postDropdownMenuAction()),
        // );
        //
        // registry.registerFileUploadMethod(
        //     <FileUploadMethodIcon/>,
        //     () => store.dispatch(fileUploadMethodAction()),
        //     'Upload using Demo Plugin',
        // );
        //
        // registry.registerWebSocketEventHandler(
        //     'custom_' + PluginId + '_status_change',
        //     (message) => {
        //         store.dispatch(websocketStatusChange(message));
        //     },
        // );

        registry.registerReducer(reducer);

        // Immediately fetch the current plugin status.
        store.dispatch(getStatus());

        // Fetch the current status whenever we recover an internet connection.
        registry.registerReconnectHandler(() => {
            store.dispatch(getStatus());
        });
    }

    uninitialize() {
        //eslint-disable-next-line no-console
        console.log(PluginId + '::uninitialize()');
    }
}
