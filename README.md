# nossl

A very simple "proxy" for stripping off SSL/TLS from a request. This is useful
for old Kindle browsers where the SSL certificates are all out of date and they
cannot connect to modern websites. It is flaky and doesn't handle everything
perfectly, but it does allow access to many more websites than the Kindle can
currently access.