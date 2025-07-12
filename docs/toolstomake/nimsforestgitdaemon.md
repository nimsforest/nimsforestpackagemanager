# NimsForest Git Daemon

A simple Makefile-based tool to serve Git repositories using git daemon.

## Makefile: `makefile.nimsforestgitdaemon`

```makefile
# makefile.nimsforestgitdaemon

help:
	@echo "NimsForest Git Daemon - Simple Git repository server"
	@echo ""
	@echo "Commands:"
	@echo "  help     - Show this help message"
	@echo "  hello    - Check dependencies (git installed, .env file exists)"
	@echo "  init     - Initialize by setting repository folder path"
	@echo "  start    - Start the git daemon server"
	@echo "  status   - Check if daemon is running and show connection info"
	@echo "  stop     - Stop the git daemon server"
	@echo "  restart  - Stop and start the daemon"
	@echo ""
	@echo "Usage:"
	@echo "  1. make -f makefile.nimsforestgitdaemon hello"
	@echo "  2. make -f makefile.nimsforestgitdaemon init"
	@echo "  3. make -f makefile.nimsforestgitdaemon start"
	@echo ""
	@echo "Then clone repositories with:"
	@echo "  git clone git://your-server-ip:9418/repo-name.git"

hello:
	@echo "Checking dependencies..."
	@which git > /dev/null || (echo "Git not installed" && exit 1)
	@test -f .env || (echo ".env file missing" && exit 1)
	@echo "Ready to go!"

init:
	@read -p "Repository folder path: " repo_path; \
	echo "REPO_BASE_PATH=$$repo_path" > .env; \
	echo "GIT_DAEMON_PORT=9418" >> .env
	@echo "Initialized with repository path"

start:
	@source .env && \
	if pgrep -f "git daemon.*$$REPO_BASE_PATH" > /dev/null; then \
		echo "Git daemon already running"; \
	else \
		echo "Starting git daemon on port $$GIT_DAEMON_PORT..."; \
		git daemon --base-path=$$REPO_BASE_PATH --export-all --reuseaddr --verbose --detach; \
		echo "Git daemon started"; \
	fi

status:
	@source .env 2>/dev/null || (echo "Not initialized. Run 'make init' first" && exit 1); \
	if pgrep -f "git daemon.*$$REPO_BASE_PATH" > /dev/null; then \
		echo "Git daemon is running"; \
		echo "Serving repositories from: $$REPO_BASE_PATH"; \
		echo "Access via: git://$(shell hostname -I | awk '{print $$1}'):$$GIT_DAEMON_PORT/repo-name.git"; \
	else \
		echo "Git daemon is not running"; \
	fi

stop:
	@source .env 2>/dev/null || (echo "Not initialized" && exit 1); \
	if pgrep -f "git daemon.*$$REPO_BASE_PATH" > /dev/null; then \
		pkill -f "git daemon.*$$REPO_BASE_PATH"; \
		echo "Git daemon stopped"; \
	else \
		echo "Git daemon is not running"; \
	fi

restart: stop start

.PHONY: help hello init start status stop restart
```

## Usage

### Initial Setup
1. Save the makefile as `makefile.nimsforestgitdaemon`
2. Check dependencies: `make -f makefile.nimsforestgitdaemon hello`
3. Initialize: `make -f makefile.nimsforestgitdaemon init`
4. Start the daemon: `make -f makefile.nimsforestgitdaemon start`

### Managing the Daemon
- **Check status**: `make -f makefile.nimsforestgitdaemon status`
- **Stop daemon**: `make -f makefile.nimsforestgitdaemon stop`  
- **Restart daemon**: `make -f makefile.nimsforestgitdaemon restart`

### Accessing Repositories
Once the daemon is running, clone repositories with:
```bash
git clone git://your-server-ip:9418/repo-name.git
```

## Features
- **Single daemon** serves multiple repositories
- **Automatic discovery** of Git repositories in the configured folder
- **Read-only access** via git:// protocol
- **VLAN-friendly** for trusted network environments
- **Minimal dependencies** - just Git and standard Unix tools

## File Structure
The tool creates a `.env` file with:
```
REPO_BASE_PATH=/path/to/your/repositories
GIT_DAEMON_PORT=9418
```

## Repository Layout
Your repository folder should contain bare Git repositories:
```
/your-repo-folder/
├── project-a.git/
├── project-b.git/
└── project-c.git/
```

Each can be accessed as:
- `git://server-ip:9418/project-a.git`
- `git://server-ip:9418/project-b.git`
- `git://server-ip:9418/project-c.git`