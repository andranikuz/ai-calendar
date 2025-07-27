#!/bin/bash

# Script to update all antd imports to use centralized import file

# Find all TypeScript/React files
FILES=$(find src -name "*.tsx" -o -name "*.ts" | grep -v "utils/antd.ts")

for file in $FILES; do
  # Check if file contains antd imports
  if grep -q "from 'antd'" "$file"; then
    echo "Updating $file..."
    
    # Get relative path from file to utils/antd.ts
    DIR=$(dirname "$file")
    RELATIVE_PATH=$(python3 -c "import os.path; print(os.path.relpath('src/utils/antd', '$DIR'))")
    
    # Replace the import statement
    sed -i.bak "s|from 'antd'|from '$RELATIVE_PATH'|g" "$file"
    
    # Remove backup file
    rm "${file}.bak"
  fi
done

echo "All antd imports have been updated!"