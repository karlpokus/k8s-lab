# app
- gw: the gateway. Requires auth via the user service. Calls blog service if successful.
- user: r/w the users db
- blog: r/w the blog db

# todos
- [x] GET /posts
- [x] GET /post/:title
- [x] POST /post
- [ ] api tests
- [ ] run test on merge
- [ ] readiness and liveness probes
- [ ] add binary version
- [ ] add graceful exits
- [x] add request logs toggle from env
- [ ] create configmap from file

# usage
We'll use docker-compose to orchestrate containers for local development and we'll keep images simple i.e huge - not optimized, to cut down on rebuild time.

```bash
# run all
$ docker-compose -f mongo.yml -f apps.yml up &
# rebuild image and restart container on src file change
$ watchexec -w src/gw "docker-compose up -d --build gw" & \
watchexec -w src/user "docker-compose up -d --build user" & \
watchexec -w src/blog "docker-compose up -d --build blog" &
```

known bugs
- docker-compose service names change colour between restarts
- go build errors end up in docker process list

# deploy
docker tag <latest> <tag>
docker push <img>
kubectl apply -f app/src -R && kubectl get pods -w -o wide
note: if only configmap change: kubectl rollout restart deploy/<thing>
