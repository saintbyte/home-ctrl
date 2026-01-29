# System Service Examples

This directory contains example configuration files for running home-ctrl as a system service.

## Available Configurations

### systemd (Modern Linux systems)
- **File**: `home-ctrl.service`
- **Location**: `/etc/systemd/system/home-ctrl.service`
- **Commands**:
  ```bash
  # Copy the service file
  sudo cp examples/system/home-ctrl.service /etc/systemd/system/
  
  # Reload systemd
  sudo systemctl daemon-reload
  
  # Enable and start the service
  sudo systemctl enable home-ctrl
  sudo systemctl start home-ctrl
  
  # Check status
  sudo systemctl status home-ctrl
  ```

### SystemV init (Older Linux systems)
- **File**: `home-ctrl.init`
- **Location**: `/etc/init.d/home-ctrl`
- **Commands**:
  ```bash
  # Copy the init script
  sudo cp examples/system/home-ctrl.init /etc/init.d/home-ctrl
  
  # Make it executable
  sudo chmod +x /etc/init.d/home-ctrl
  
  # Add to startup (on RedHat/CentOS)
  sudo chkconfig --add home-ctrl
  sudo chkconfig home-ctrl on
  
  # Start the service
  sudo service home-ctrl start
  
  # Check status
  sudo service home-ctrl status
  ```

## Configuration Notes

1. **Paths**: Update the paths in the configuration files to match your installation:
   - `ExecStart` / `DAEMON`: Path to your home-ctrl binary
   - `WorkingDirectory`: Working directory for the service
   - `Environment`: Configuration file path

2. **User**: It's recommended to create a dedicated user for the service:
   ```bash
   sudo useradd -r -s /bin/false home-ctrl
   ```

3. **Permissions**: Ensure the service user has access to required resources.

4. **Logging**: Consider adding logging configuration if needed.