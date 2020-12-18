# rubocop:disable all

require 'strscan'

def parse(input)
  ss = StringScanner.new(input.gsub(" ", ""))

  res = []

  loop do
    number = ss.scan(/[0-9]+/)
    if number
      res << number.to_i
      next
    end

    op = ss.scan(/[+*]/)
    if op
      res << op
      next
    end

    open = ss.scan(/\(/)
    if open
      res << "("
      next
    end

    close = ss.scan(/\)/)
    if close
      res << ")"
      next
    end

    break
  end

  res
end

def find_close(tokens)
  depth = 0
  res = -1

  tokens.each.with_index do |e, i|
    if e == ")" && depth == 1
      res = i
      break
    end

    if e == ")"
      depth -= 1
    end

    if e == "("
      depth += 1
    end
  end

  res
end

def to_ast(tokens)
  res = []

  while tokens.size > 0
    if tokens[0] == "("
      s = 1
      e = find_close(tokens)

      res << to_ast(tokens[s...e])
      tokens = tokens.drop(e+1)
      next
    end

    res << tokens.shift
  end

  res
end

def eval_plus(tokens)
  res = []

  while tokens.size >= 1
    left = tokens.shift
    op = tokens.shift

    case op
    when "+"
      right = tokens.shift
      tokens.unshift(left + right)
    when nil
      res << left
    else
      res << left
      res << op
    end
  end

  res
end

def eval_times(tokens)
  res = []

  while tokens.size >= 1
    left = tokens.shift
    op = tokens.shift

    case op
    when "*"
      right = tokens.shift
      tokens.unshift(left * right)
    when nil
      res << left
    else
      res << left
      res << op
    end
  end

  res
end

def calc(a)
  return a if a.class != Array

  if a.class == Array
    a.map! { |e| calc(e) }
  end

  a = eval_plus(a)
  a = eval_times(a)

  a[0]
end

res = File.read(ARGV[0]).split("\n").map do |line|
  calc(to_ast(parse(line)))
end.inject(:+)

puts res
