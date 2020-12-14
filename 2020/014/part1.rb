# rubocop:disable all

input = File.read(ARGV[0]).split("\n")

mask_base = nil
mask_val = nil

mem = {}

def as_bits_str(value)
  (value | (1 << 36)).to_s(2)
end

def mask_it(value, mask_base, mask_val)
  puts as_bits_str(value)
  puts as_bits_str(mask_base)
  puts as_bits_str(mask_val)
  puts "-----"
  puts as_bits_str(value & mask_base)
  puts as_bits_str(value & mask_base | mask_val)


  (value & mask_base) | mask_val
end

input.each do |line|
  left, right = line.split(" = ")

  if left == "mask"
    mask_base = right.split("").map { |c| c == "X" ? 1 : 0 }.join("").to_i(2)
    mask_val  = right.gsub("X", "0").to_i(2)

    puts "Changing mask to #{mask_base} #{mask_val}"
    next
  end

  location = left.gsub("[", "").gsub("]", "").gsub("mem", "").to_i
  value = mask_it(right.to_i, mask_base, mask_val)

  puts "Setting #{location} to #{value}"
  mem[location] = value
end

puts mem.map { |k, v| v }.inject(:+)
