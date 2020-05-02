# app
- gw: the gateway. Requires auth via the user service. Calls blog service if successful.
- user: r/w the users db
- blog: r/w the blog db

# todos
- [x] api /posts
- [x] api /post/:title
- [ ] api tests
- [ ] run test on merge
- [ ] readiness and liveness probes

# usage
We'll use docker-compose to orchestrate containers for local development and we'll keep images simple i.e huge - not optimized, to cut down on rebuild time.

```bash
# rebuild image and restart container on src file change
$ watchexec -w src/gw "docker-compose up -d --build gw" & \
 watchexec -w src/user "docker-compose up -d --build user" & \
 watchexec -w src/blog "docker-compose up -d --build blog" &
# run all
$ docker-compose up [--build]
```

known bugs
- docker-compose services change name between restarts
- go build errors end up in docker process list

# deploy
docker build <dockerfile> -t <tag> (or re-tag latest)
docker push <img>
scp <app>/deployment.yaml to host
kubectl apply -f <app>/deployment.yaml
