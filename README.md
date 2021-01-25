# PaySuper MongoDB Driver

[![License: GPL 3.0](https://img.shields.io/badge/License-GPL3.0-green.svg)](https://opensource.org/licenses/Gpl3.0)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/paysuper/paysuper-database-mongo/issues)
[![Build Status](https://github.com/paysuper/paysuper-database-mongo/workflows/Build/badge.svg?branch=master)](https://github.com/paysuper/paysuper-database-mongo/actions) 
[![codecov](https://codecov.io/gh/paysuper/paysuper-database-mongo/branch/master/graph/badge.svg)](https://codecov.io/gh/paysuper/paysuper-database-mongo)
[![go report](https://goreportcard.com/badge/github.com/paysuper/paysuper-database-mongo)](https://goreportcard.com/report/github.com/paysuper/paysuper-database-mongo)

PaySuper MongoDB Driver is a Mongo MGO library wrapper.

***

## Table of Contents

- [Usage](#usage)
- [Contributing](#contributing-feature-requests-and-support)
- [License](#license)

# Usage

Application handles configurations from the environment variables.

### Environment variables:

| Name               | Required | Default  | Description                     |
|:-------------------|:--------:|:---------|:--------------------------------|
| `MONGO_DIAL_TIMEOUT` | -        | 10       | MongoDB dial timeout in seconds |
| `MONGO_DSN`          | true     | -        | MongoDB DSN connection string   |

## Contributing, Feature Requests and Support

If you like this project then you can put a ⭐ on it. It means a lot to us.

If you have an idea of how to improve PaySuper (or any of the product parts) or have general feedback, you're welcome to submit a [feature request](../../issues/new?assignees=&labels=&template=feature_request.md&title=).

Chances are, you like what we have already but you may require a custom integration, a special license or something else big and specific to your needs. We're generally open to such conversations.

If you have a question and can't find the answer yourself, you can [raise an issue](../../issues/new?assignees=&labels=&template=issue--support-request.md&title=I+have+a+question+about+<this+and+that>+%5BSupport%5D) and describe what exactly you're trying to do. We'll do our best to reply in a meaningful time.

We feel that a welcoming community is important and we ask that you follow PaySuper's [Open Source Code of Conduct](https://github.com/paysuper/code-of-conduct/blob/master/README.md) in all interactions with the community.

PaySuper welcomes contributions from anyone and everyone. Please refer to [our contribution guide to learn more](CONTRIBUTING.md).

## License

The project is available as open source under the terms of the [GPL v3 License](https://www.gnu.org/licenses/gpl-3.0).
