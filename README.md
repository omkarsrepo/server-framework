## Golang GIN Production Server Framework

### Features
1. Router
2. Configuration
3. Middlewares
4. Shutdown Hook
5. Custom Exception
6. Custom Json Validators and Validators Exception
7. Json Utils
8. Secret Service
9. Database Initialization
10. In Memory Cache Service
11. Logger Service
12. Commands
13. Rate Limiter Middleware
14. Request Timeout Middleware
15. Request Logger Middleware 
16. TraceId support for Logger and Errors 
17. Auto Go Max Processes Configuration 
18. Go Memory Limit Configuration 
19. Graceful Shutdown
20. Default Cors Configuration with override support
21. Enabled PProf
```
curl -v -H "Authorization: foobar" -o profile.pb.gz \
  http://localhost:8080/metrics/pprof/profile?seconds=60
go tool pprof -http=:8099 profile.pb.gz
```
22. Enabled Gzip Compression with option to exclude paths.
23. Health Ping Endpoint: `/health/IhEaf/ping`
