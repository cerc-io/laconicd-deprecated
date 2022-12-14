local config = import 'default.jsonnet';

config {
  'laconic_9000-1'+: {
    validators: super.validators + [{
      name: 'fullnode',
    }],
  },
}
