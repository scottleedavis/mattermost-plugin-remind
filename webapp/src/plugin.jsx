import React from 'react';

import {FormattedMessage} from 'react-intl';

import en from 'i18n/en.json';

import es from 'i18n/es.json';

import {id as pluginId} from './manifest';

import Root from './components/root';
// import BottomTeamSidebar from './components/bottom_team_sidebar';
import LeftSidebarHeader from './components/left_sidebar_header';
import SubMenu from './components/sub_menu';
// import LinkTooltip from './components/link_tooltip';
// import UserAttributes from './components/user_attributes';
// import UserActions from './components/user_actions';
// import PostType from './components/post_type';
// import {
//     MainMenuMobileIcon,
//     ChannelHeaderButtonIcon,
//     FileUploadMethodIcon,
// } from './components/icons';
import {
    mainMenuAction,
    channelHeaderButtonAction,
    fileUploadMethodAction,
    postDropdownMenuAction,
    websocketStatusChange,
    getStatus,
} from './actions';
import reducer from './reducer';

function getTranslations(locale) {
    switch (locale) {
    case 'en':
        return en;
    case 'es':
        return es;
    }
    return {};
}

export default class DemoPlugin {
    initialize(registry, store) {
        registry.registerRootComponent(Root);
        // registry.registerPopoverUserAttributesComponent(UserAttributes);
        // registry.registerPopoverUserActionsComponent(UserActions);
        // registry.registerLeftSidebarHeaderComponent(LeftSidebarHeader);
        // registry.registerLinkTooltipComponent(LinkTooltip);
        // registry.registerBottomTeamSidebarComponent(
        //     BottomTeamSidebar,
        // );

        // registry.registerChannelHeaderButtonAction(
        //     <ChannelHeaderButtonIcon/>,
        //     () => store.dispatch(channelHeaderButtonAction()),
        //     <FormattedMessage
        //         id='plugin.name'
        //         defaultMessage='Demo Plugin'
        //     />,
        // );
        //
        // registry.registerPostTypeComponent('custom_demo_plugin', PostType);
        //
        // registry.registerMainMenuAction(
        //     <FormattedMessage
        //         id='plugin.name'
        //         defaultMessage='Demo Plugin'
        //     />,
        //     () => store.dispatch(mainMenuAction()),
        //     <MainMenuMobileIcon/>,
        // );

        registry.registerPostDropdownMenuAction(SubMenu);

        registry.registerPostDropdownMenuAction(
            <FormattedMessage
                id='plugin.name'
                defaultMessage='Demo Plugin'
            />,
            () => store.dispatch(postDropdownMenuAction()),
        );

        // registry.registerFileUploadMethod(
        //     <FileUploadMethodIcon/>,
        //     () => store.dispatch(fileUploadMethodAction()),
        //     <FormattedMessage
        //         id='plugin.upload'
        //         defaultMessage='Upload using Demo Plugin'
        //     />,
        // );

        // registry.registerWebSocketEventHandler(
        //     'custom_' + pluginId + '_status_change',
        //     (message) => {
        //         store.dispatch(websocketStatusChange(message));
        //     },
        // );

        registry.registerReducer(reducer);

        // // Immediately fetch the current plugin status.
        // store.dispatch(getStatus());
        //
        // // Fetch the current status whenever we recover an internet connection.
        // registry.registerReconnectHandler(() => {
        //     store.dispatch(getStatus());
        // });

        registry.registerTranslations(getTranslations);
    }

    uninitialize() {
        //eslint-disable-next-line no-console
        console.log(pluginId + '::uninitialize()');
    }
}
