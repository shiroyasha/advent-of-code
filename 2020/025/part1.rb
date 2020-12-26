# card = 5764801
# door = 17807724

card = 15733400
door = 6408062

def secret_loop(input)
  n = 1
  res = 1

  loop do
    res = res * 7 % 20201227

    if res == input
      return n
    end

    n += 1
  end
end

a = secret_loop(card)
b = secret_loop(door)

res = 1

a.times do
  res = (res * door) % 20201227
end

puts res
