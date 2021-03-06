
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

SUFFIXES = .go .c
OBJEXT= none

EXTRA_PROGRAMS = builder
builder_SOURCES = builder.go

builder$(EXEEXT): $(builder_SOURCES)
	$(GO) build -o builder$(EXEEXT) $(builder_SOURCES)

sbin_PROGRAMS = m5app m5pwd

m5app_SOURCES = m5app.go
m5pwd_SOURCES = m5pwd.go

GOFLAGS = -ldflags="-s -w"

EXTRA_m5app_SOURCES = \
	server/server.go \
	tools/file.go \
	server/config/config.go \
	server/daemon/daemon.go \
	server/controller/hello/hello.go \
	server/middleware/requestlog.go \
	server/middleware/responselog.go \
	model/user/user.go

EXTRA_m5app_SOURCES += \
	server/bundle/public.go

EXTRA_m5pwd_SOURCES = \
	server/config/config.go \
	model/user/user.go

server/bundle/public.go: builder$(EXEEXT) public/index.html
	./builder$(EXEEXT) --package=bundle -o server/bundle/public.go public/

m5app$(EXEEXT): $(m5app_SOURCES) $(EXTRA_m5app_SOURCES)
	$(GO) build $(GOFLAGS) -o m5app$(EXEEXT) $(m5app_SOURCES)

m5pwd$(EXEEXT): $(m5pwd_SOURCES) $(EXTRA_m5pwd_SOURCES)
	$(GO) build $(GOFLAGS) -o m5pwd$(EXEEXT) $(m5pwd_SOURCES)

EXTRA_DIST = \
	README.md \
	go.mod \
	go.sum

clean-local:
	rm -f builder$(EXEEXT)
	rm -rf autom4te.cache

install-data-local:
	test -z $(DESTDIR)$(APP_LIBDIR) || $(MKDIR_P) $(DESTDIR)$(APP_LIBDIR)
#	cp -r public $(DESTDIR)$(APP_LIBDIR)
	test -z $(DESTDIR)$(APP_CONFDIR) || $(MKDIR_P) $(DESTDIR)$(APP_CONFDIR)
	test -z $(DESTDIR)$(APP_LOGDIR) || $(MKDIR_P) $(DESTDIR)$(APP_LOGDIR)
	test -z $(DESTDIR)$(APP_RUNDIR) || $(MKDIR_P) $(DESTDIR)$(APP_RUNDIR)
	test -z $(DESTDIR)$(APP_DATABASEDIR) || $(MKDIR_P) $(DESTDIR)$(APP_DATABASEDIR)
#EOF
