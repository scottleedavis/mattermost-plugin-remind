
import en from 'i18n/en.json';

import es from 'i18n/es.json';

import {id as pluginId} from './manifest';

import RemindMenuItem from './components/remind_menu_item';

import {
    postDropdownMenuAction,
} from './actions';

// import reducer from './reducer';

function getTranslations(locale) {
    switch (locale) {
    case 'en':
        return en;
    case 'es':
        return es;
    }
    return {};
}

export default class RemindPlugin {
    initialize(registry, store) {
        registry.registerPostDropdownMenuAction(
            <RemindMenuItem/>,
            (postId) => store.dispatch(postDropdownMenuAction(postId)),
        );

        // registry.registerReducer(reducer);

        registry.registerTranslations(getTranslations);
    }

    uninitialize() {
        //eslint-disable-next-line no-console
        console.log(pluginId + '::uninitialize()');
    }
}
