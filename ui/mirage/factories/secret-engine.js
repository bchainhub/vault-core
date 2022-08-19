import { Factory } from 'ember-cli-mirage';

export default Factory.extend({
  path: 'foo/',
  description: 'secret-engine generated by mirage',
  local: true,
  sealWrap: true,
  // set in afterCreate
  accessor: 'type_7f52940',
  type: 'kv',
  options: null,

  afterCreate(secretEngine) {
    if (!secretEngine.options && ['generic', 'kv'].includes(secretEngine.type)) {
      secretEngine.options = {
        version: '2',
      };
    }
  },
});
