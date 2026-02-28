import os
import sys

class FileManager:
    @staticmethod
    def safe_open(file_path, mode):
        try:
            return open(file_path, mode)
        except IOError as e:
            print(f"Error opening file {file_path}: {e}")
            sys.exit(1)

class CommandExecutor:
    @staticmethod
    def safe_execute(command):
        try:
            os.system(command)
        except Exception as e:
            print(f"Error executing command '{command}': {e}")
            sys.exit(1)


def generate_config_and_run(config_path):
    try:
        # Load or generate configuration
        with FileManager.safe_open(config_path, 'r') as file:
            config = file.read()
            # Process configuration
        print("Configuration successful")
        CommandExecutor.safe_execute('run_application')
    except Exception as e:
        print(f"An error occurred: {e}")
        sys.exit(1)

# Example usage
if __name__ == '__main__':
    generate_config_and_run('config.json')
