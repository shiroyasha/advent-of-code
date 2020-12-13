# rubocop:disable all

def extended_gcd(a, b)
  last_remainder, remainder = a.abs, b.abs

  x, last_x, y, last_y = 0, 1, 1, 0

  while remainder != 0
    last_remainder, (quotient, remainder) = remainder, last_remainder.divmod(remainder)
    x, last_x = last_x - quotient*x, x
    y, last_y = last_y - quotient*y, y
  end

  return last_remainder, last_x * (a < 0 ? -1 : 1)
end

def invmod(e, et)
  g, x = extended_gcd(e, et)

  if g != 1
    raise 'Multiplicative inverse modulo does not exist!'
  end

  x % et
end

def chinese_remainder(mods, remainders)
  max = mods.inject(:*)

  series = remainders.zip(mods).map {|r,m| (r * max * invmod(max/m, m) / m) }

  series.inject(:+) % max
end

def solution(input)
  mods = []
  rems = []

  input.split(",").each.with_index do |n, i|
    next if n == "x"

    mods << n.to_i
    rems << -i
  end

  puts "#{mods.inspect} #{rems.inspect}"

  chinese_remainder(mods, rems)
end

puts solution("7,13,x,x,59,x,31,19")
puts solution("17,x,13,19")
puts solution("67,7,59,61")
puts solution("67,x,7,59,61")
puts solution("67,7,x,59,61")
puts solution("1789,37,47,1889")
puts solution("19,x,x,x,x,x,x,x,x,41,x,x,x,37,x,x,x,x,x,367,x,x,x,x,x,x,x,x,x,x,x,x,13,x,x,x,17,x,x,x,x,x,x,x,x,x,x,x,29,x,373,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,23")
