import os

def generate_large_file(filename, size_gb):
    size_bytes = size_gb * 1024 * 1024 * 1024
    chunk_size = 1024 * 1024  # 1MB chunks
    
    with open(filename, 'wb') as f:
        remaining = size_bytes
        while remaining > 0:
            chunk = min(chunk_size, remaining)
            # Use os.urandom for random data, or b'\x00' * chunk for zeros
            data = os.urandom(chunk)  # or b'\x00' * chunk for faster generation
            f.write(data)
            remaining -= chunk
            
            # Progress indicator
            if (size_bytes - remaining) % (100 * 1024 * 1024) == 0:
                print(f"Written: {(size_bytes - remaining) // (1024**3):.1f} GB")

generate_large_file("test_40gb.bin", 5)