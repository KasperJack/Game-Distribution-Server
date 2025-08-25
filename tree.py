#!/usr/bin/env python3
"""
Network Directory Tree Generator
Generates directory structure optimized for network transfer:
1. Directories first (so client can create folder structure)
2. Files second (so client can request them)
"""

import os
import json
import argparse
from pathlib import Path
from datetime import datetime

def scan_directory_for_network(directory, show_hidden=False, include_stats=False):
    """
    Scan directory and return data optimized for network transfer.
    Returns directories first, then files.

    Args:
        directory (Path): Directory to scan
        show_hidden (bool): Include hidden files
        include_stats (bool): Include file statistics

    Returns:
        dict: Contains 'directories' and 'files' lists
    """
    directories = []
    files = []

    def process_item(path, relative_to, item_type):
        try:
            relative_path = path.relative_to(relative_to)

            item_data = {
                'path': str(relative_path).replace('\\', '/'),  # Normalize path separators
                'name': path.name,
                'parent': str(relative_path.parent).replace('\\', '/') if relative_path.parent != Path('.') else '',
                'depth': len(relative_path.parts) - 1
            }

            if item_type == 'file' and include_stats:
                try:
                    stat = path.stat()
                    item_data.update({
                        'size': stat.st_size,
                        'modified': datetime.fromtimestamp(stat.st_mtime).isoformat(),
                        'checksum_needed': True  # Flag for client to verify after transfer
                    })
                except (OSError, PermissionError):
                    item_data['size'] = 0
                    item_data['error'] = 'Cannot read stats'

            return item_data

        except (PermissionError, OSError) as e:
            return {
                'path': str(path.relative_to(relative_to)).replace('\\', '/'),
                'name': path.name,
                'error': str(e),
                'parent': str(path.parent.relative_to(relative_to)).replace('\\', '/') if path.parent != relative_to else '',
                'depth': len(path.relative_to(relative_to).parts) - 1
            }

    # Walk through directory
    for root, dirs, filenames in os.walk(directory):
        root_path = Path(root)

        # Filter hidden items
        if not show_hidden:
            dirs[:] = [d for d in dirs if not d.startswith('.')]
            filenames = [f for f in filenames if not f.startswith('.')]

        # Add directories (except root)
        if root_path != directory:
            dir_data = process_item(root_path, directory, 'directory')
            directories.append(dir_data)

        # Add files
        for filename in filenames:
            file_path = root_path / filename
            file_data = process_item(file_path, directory, 'file')
            files.append(file_data)

    # Sort directories by depth then path (ensures parent dirs come before children)
    directories.sort(key=lambda x: (x['depth'], x['path']))

    # Sort files by path
    files.sort(key=lambda x: x['path'])

    return {
        'directories': directories,
        'files': files
    }

def save_network_format(data, output_file, root_path, format_type='json'):
    """Save in network-optimized format"""

    output = {
        'protocol_version': '1.0',
        'root_name': root_path.name,
        'root_path': str(root_path),
        'generated': datetime.now().isoformat(),
        'total_directories': len(data['directories']),
        'total_files': len(data['files']),
        'transfer_order': {
            'step_1': 'Create all directories in order',
            'step_2': 'Request files by path'
        },
        'directories': data['directories'],
        'files': data['files']
    }

    if format_type == 'json':
        with open(output_file, 'w', encoding='utf-8') as f:
            json.dump(output, f, indent=2, ensure_ascii=False)

    elif format_type == 'protocol':
        # Custom protocol format: easy to parse line by line
        with open(output_file, 'w', encoding='utf-8') as f:
            f.write(f"PROTOCOL_VERSION:1.0\n")
            f.write(f"ROOT_NAME:{root_path.name}\n")
            f.write(f"GENERATED:{datetime.now().isoformat()}\n")
            f.write(f"TOTAL_DIRS:{len(data['directories'])}\n")
            f.write(f"TOTAL_FILES:{len(data['files'])}\n")
            f.write("BEGIN_DIRECTORIES\n")

            for dir_item in data['directories']:
                line = f"DIR:{dir_item['path']}:{dir_item['depth']}"
                if 'error' in dir_item:
                    line += f":ERROR:{dir_item['error']}"
                f.write(line + '\n')

            f.write("BEGIN_FILES\n")

            for file_item in data['files']:
                line = f"FILE:{file_item['path']}"
                if 'size' in file_item:
                    line += f":{file_item['size']}"
                if 'error' in file_item:
                    line += f":ERROR:{file_item['error']}"
                f.write(line + '\n')

            f.write("END_MANIFEST\n")

def generate_client_example(output_dir):
    """Generate example client code"""
    client_code = '''#!/usr/bin/env python3
"""
Example client code for processing the directory manifest
"""

import json
import os
from pathlib import Path

def process_manifest(manifest_file, target_dir):
    """Process manifest and create directory structure"""

    with open(manifest_file, 'r') as f:
        if manifest_file.endswith('.json'):
            data = json.load(f)
        else:
            # Parse protocol format
            data = parse_protocol_format(f.read())

    target_path = Path(target_dir)
    target_path.mkdir(exist_ok=True)

    print(f"Creating directory structure in: {target_path}")

    # Step 1: Create all directories
    print("Step 1: Creating directories...")
    for dir_item in data['directories']:
        if 'error' not in dir_item:
            dir_path = target_path / dir_item['path']
            dir_path.mkdir(parents=True, exist_ok=True)
            print(f"  Created: {dir_path}")

    # Step 2: List files to request
    print("\\nStep 2: Files to request:")
    file_requests = []
    for file_item in data['files']:
        if 'error' not in file_item:
            file_requests.append({
                'remote_path': file_item['path'],
                'local_path': target_path / file_item['path'],
                'size': file_item.get('size', 0)
            })
            print(f"  Request: {file_item['path']} ({file_item.get('size', '?')} bytes)")

    return file_requests

def parse_protocol_format(content):
    """Parse custom protocol format"""
    lines = content.strip().split('\\n')
    data = {'directories': [], 'files': []}

    section = None
    for line in lines:
        if line.startswith('BEGIN_DIRECTORIES'):
            section = 'directories'
        elif line.startswith('BEGIN_FILES'):
            section = 'files'
        elif line.startswith('DIR:') and section == 'directories':
            parts = line[4:].split(':')
            data['directories'].append({
                'path': parts[0],
                'depth': int(parts[1]) if len(parts) > 1 else 0
            })
        elif line.startswith('FILE:') and section == 'files':
            parts = line[5:].split(':')
            item = {'path': parts[0]}
            if len(parts) > 1 and parts[1].isdigit():
                item['size'] = int(parts[1])
            data['files'].append(item)

    return data

if __name__ == "__main__":
    import sys
    if len(sys.argv) != 3:
        print("Usage: python client_example.py manifest_file target_directory")
        sys.exit(1)

    file_requests = process_manifest(sys.argv[1], sys.argv[2])
    print(f"\\nReady to request {len(file_requests)} files!")
'''

    with open(f"{output_dir}/client_example.py", 'w') as f:
        f.write(client_code)

def main():
    parser = argparse.ArgumentParser(description="Generate directory tree for network transfer")
    parser.add_argument("directory", nargs="?", default=".", help="Directory to scan")
    parser.add_argument("-o", "--output", default="manifest", help="Output file name (no extension)")
    parser.add_argument("-f", "--format", choices=['json', 'protocol'], default='json',
                       help="Output format (default: json)")
    parser.add_argument("-a", "--all", action="store_true", help="Include hidden files")
    parser.add_argument("-s", "--stats", action="store_true", help="Include file statistics")
    parser.add_argument("--client-example", action="store_true", help="Generate example client code")

    args = parser.parse_args()

    # Setup paths
    target_dir = Path(args.directory).resolve()
    output_file = f"{args.output}.{args.format}"

    # Validate directory
    if not target_dir.exists() or not target_dir.is_dir():
        print(f"Error: '{target_dir}' is not a valid directory.")
        return

    print(f"Scanning directory: {target_dir}")
    print(f"Output format: {args.format}")
    print(f"Output file: {output_file}")

    # Scan directory
    data = scan_directory_for_network(target_dir, args.all, args.stats)

    # Save manifest
    save_network_format(data, output_file, target_dir, args.format)

    print(f"\\nManifest generated successfully!")
    print(f"- Directories: {len(data['directories'])} (will be created first)")
    print(f"- Files: {len(data['files'])} (will be requested after)")
    print(f"- Saved to: {output_file}")

    # Generate client example
    if args.client_example:
        generate_client_example(".")
        print(f"- Client example: client_example.py")

    # Show first few items for verification
    if data['directories']:
        print("\\nFirst directories to create:")
        for i, dir_item in enumerate(data['directories'][:5]):
            print(f"  {i+1}. {dir_item['path']} (depth: {dir_item['depth']})")
        if len(data['directories']) > 5:
            print(f"  ... and {len(data['directories']) - 5} more")

    if data['files']:
        print("\\nFirst files to request:")
        for i, file_item in enumerate(data['files'][:5]):
            size_info = f" ({file_item['size']} bytes)" if 'size' in file_item else ""
            print(f"  {i+1}. {file_item['path']}{size_info}")
        if len(data['files']) > 5:
            print(f"  ... and {len(data['files']) - 5} more")

if __name__ == "__main__":
    main()
