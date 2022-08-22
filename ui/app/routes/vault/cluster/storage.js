import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';

export default Route.extend(ClusterRoute, {
  model() {
    // findAll method will return all records in store as well as response from server
    // when removing a peer via the cli, stale records would continue to appear until refresh
    // query method will only return records from response
    return this.store.query('server', {});
  },

  actions: {
    doRefresh() {
      this.refresh();
    },
  },
});
