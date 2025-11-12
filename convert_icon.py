#!/usr/bin/env python3
"""
Convert PNG to ICO format for Windows application icon
"""

from PIL import Image
import os

def png_to_ico(png_path, ico_path):
    """Convert PNG image to ICO format with multiple sizes"""
    try:
        # Open the PNG image
        img = Image.open(png_path)
        
        # Convert to RGBA if necessary
        if img.mode != 'RGBA':
            img = img.convert('RGBA')
        
        # Create different sizes for ICO (Windows standard)
        sizes = [(256, 256), (128, 128), (64, 64), (48, 48), (32, 32), (16, 16)]
        
        # Resize images
        icons = []
        for size in sizes:
            resized = img.resize(size, Image.Resampling.LANCZOS)
            icons.append(resized)
        
        # Save