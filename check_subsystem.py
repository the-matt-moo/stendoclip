import struct
import sys

with open('stendoclip.exe', 'rb') as f:
    data = f.read(256)

pe_offset = struct.unpack_from('<I', data, 0x3C)[0]
print(f'PE header offset: 0x{pe_offset:04X}')
sig = data[pe_offset:pe_offset+4]
print(f'PE signature: {sig}')
machine = struct.unpack_from('<H', data, pe_offset+4)[0]
num_sections = struct.unpack_from('<H', data, pe_offset+6)[0]
size_opt_hdr = struct.unpack_from('<H', data, pe_offset+20)[0]
print(f'Machine: 0x{machine:04X} (0x8664=AMD64, 0x014C=x86)')
print(f'NumberOfSections: {num_sections}')
print(f'SizeOfOptionalHeader: {size_opt_hdr}')
opt_hdr_start = pe_offset + 24
subsystem = struct.unpack_from('<H', data, opt_hdr_start + 68)[0]
print(f'Subsystem: {subsystem} (2=GUI, 3=Console)')
