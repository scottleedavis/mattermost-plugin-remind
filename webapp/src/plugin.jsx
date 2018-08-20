import React from 'react';

import PluginId from './plugin_id';

import Root from './components/root';

import {
    ChannelHeaderButtonIcon,
} from './components/icons';
import {
    channelHeaderButtonAction,
    getStatus,
} from './actions';
import reducer from './reducer';

export default class DemoPlugin {
    initialize(registry, store) {
        registry.registerRootComponent(Root);

        registry.registerChannelHeaderButtonAction(
            <ChannelHeaderButtonIcon/>,
            () => store.dispatch(channelHeaderButtonAction()),
            'Remind Plugin',
        );

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
