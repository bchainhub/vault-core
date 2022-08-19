import { inject as service } from '@ember/service';
import { computed } from '@ember/object';
import Controller, { inject as controller } from '@ember/controller';
import utils from 'vault/lib/key-utils';
import ListController from 'core/mixins/list-controller';

export default Controller.extend(ListController, {
  flashMessages: service(),
  store: service(),
  clusterController: controller('vault.cluster'),

  backendCrumb: computed('clusterController.model.name', function () {
    return {
      label: 'leases',
      text: 'leases',
      path: 'vault.cluster.access.leases.list-root',
      model: this.clusterController.model.name,
    };
  }),

  isLoading: false,

  filterIsFolder: computed('filter', function () {
    return !!utils.keyIsFolder(this.filter);
  }),

  emptyTitle: computed('baseKey.id', 'filter', 'filterIsFolder', function () {
    let id = this.baseKey.id;
    let filter = this.filter;
    if (id === '') {
      return 'There are currently no leases.';
    }
    if (this.filterIsFolder) {
      if (filter === id) {
        return `There are no leases under &quot;${filter}&quot;.`;
      } else {
        return `We couldn't find a prefix matching &quot;${filter}&quot;.`;
      }
    }
    return '';
  }),

  actions: {
    revokePrefix(prefix, isForce) {
      const adapter = this.store.adapterFor('lease');
      const method = isForce ? 'forceRevokePrefix' : 'revokePrefix';
      const fn = adapter[method];
      fn.call(adapter, prefix)
        .then(() => {
          return this.transitionToRoute('vault.cluster.access.leases.list-root').then(() => {
            this.flashMessages.success(`All of the leases under ${prefix} will be revoked.`);
          });
        })
        .catch((e) => {
          const errString = e.errors.join('.');
          this.flashMessages.danger(
            `There was an error attempting to revoke the prefix: ${prefix}. ${errString}.`
          );
        });
    },
  },
});
