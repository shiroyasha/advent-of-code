# rubocop:disable all

input = File.read(ARGV[0]).split("\n")

mask = nil

mem = {}

def join(mask, address)
  address = (address | (1 << 37)).to_s(2)[2..-1].split("")

  puts address.join("")
  puts mask.join("")

  floating = mask.zip(address).map do |e|
    if e[0] == "X"
      "X"
    elsif e[0] == "1"
      "1"
    else
      e[1]
    end
  end

  floating.join("")
end

def change(str, i, val)
  res = str.clone
  res[i] = val
  res
end

def expand(address, index = 0)
  if index == 36
    return address
  end

  if address[index] == "X"
    [
      expand(change(address, index, "0"), index + 1),
      expand(change(address, index, "1"), index + 1)
    ].flatten
  else
    expand(address, index + 1)
  end
end

input.each do |line|
  left, right = line.split(" = ")

  if left == "mask"
    mask = right.split("")
    next
  end

  location = left.gsub("[", "").gsub("]", "").gsub("mem", "").to_i
  value = right.to_i

  puts "Joining"
  joined = join(mask, location)

  puts "#{joined}"
  addresses = expand(joined)

  addresses.each do |address|
    puts "Setting #{address} #{address.to_i(2)}"
    mem[address] = value
  end
end

puts mem.map { |k, v| v }.inject(:+)
