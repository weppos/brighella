# Brighella

_Brighella_ is a simple URL-masking redirect service built on Go.

It is designed to be deployed on a PaaS service such as Heroku, and it requires zero configuration. The redirect target is fetched from a specific DNS record attached to the requested host.

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/weppos/brighella)

[Brighella](https://en.wikipedia.org/wiki/Brighella) is a comic, masked character from the Commedia dell'arte. It's also an Italian carnival mask.

## How it works

The redirect is performed using an iframe. The target URL is rendered in an iframe, so that the URL will continue to match your URL in the browser, but the user will effectively navigate in the target URL.

This is known as _iframe masking_, _masked redirect_, _masked forwarding_ or _iframe redirect_.

Please note that the target site may explicitly forbid the displaying of the content within an iframe (using the [`X-Frame-Options` header](https://developer.mozilla.org/en-US/docs/Web/HTTP/X-Frame-Options)). In that case, there is nothing you can do because that's what the owner of the target URL wants.

## Usage

Let's assume you want to redirect `example.com` to `http://somesite.com`. What you have to do is:

- Deploy the application
- Make sure the app is properly configured to respond to the `example.com` domain (for Heroku use the `domains:add example.com` command)
- Make sure the DNS record for `example.com` points to the server where the app is deployed
- Configure a DNS TXT record called `_frame.example.com` with the content `http://somesite.com`

That's it. The app will automatically try to load the target of the redirect from the DNS record.

You can configure as many domains you want for a single instance of the application, as long as each domain has a corresponding `_frame` record.

Important: that the record name MUST match the domain name, prefixed with `_frame`. Here's some examples of expected configurations:

```
example.com requires _frame.example.com
www.example.com requires _frame.www.example.com
subdomain.example.com requires _frame.subdomain.example.com
```

If you need a reliable DNS provider to manage your domain, [check out DNSimple](https://dnsimple.com/). FYI, this project was born in response to the various customers at DNSimple that asked us information about how to configure an URL-masked redirect.


## License

Copyright (c) 2016 Simone Carletti. This is Free Software distributed under the MIT license.
