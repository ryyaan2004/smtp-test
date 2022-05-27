# Simple SMTP Test
A script put together in about an hour to generate email load

## Notes
* No security, this was created to test configurations on a new smtp server
* When connection errors occur it does not recover; gracefully or otherwise

## Future
I probably won't do any more work on this as I'm not likely to be in a situation where I need to quickly test an email server but aren't able to find a more suitable alternative. However, if I do get a chance here are some things I think it would be fun to address:

1. Handle login - take password from stdin
2. Allow the use of a test message template
3. Change from a simple script to a usable lib