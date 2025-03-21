
# drone-github-clone

use 443 port
```yml
  - name: clone
    image: wwma/drone-github-clone
    settings:
      SSH_KEY:
        from_secret: ssh_key
```

## Authors

- [@zzfn](https://github.com/zzfn)


## Badges

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)
[![Build Status](https://drone.zzfzzf.com/api/badges/drone-plugin/github/status.svg)](https://drone.zzfzzf.com/drone-plugin/github)
