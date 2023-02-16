# Token

> Token验证


-------------
```markdown
-----------------------------------------------------
 client             server             token storage
  auth    --->       get token    --->      token
                     save         <---      token payload
  ok?     <---       verify
next time
  auth    --->       check
  ok>     <---       verify
```
