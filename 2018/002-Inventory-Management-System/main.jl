using StatsBase

function find_checksum(path)
  file  = open(path)
  lines = readlines(file)

  count2 = 0
  count3 = 0

  for l in lines
    m = values(countmap([c for c in l]))

    if in(2, m)
      count2 = count2 + 1
    end

    if in(3, m)
      count3 = count3 + 1
    end
  end

  count2 * count3
end

res1 = find_checksum("input2.txt")
println(res1)

function diff(line1, line2)
  c = 0
  res = ""

  for i in 1:length(line1)
    if line1[i] != line2[i]
      c = c + 1
    else
      res = res * line1[i]
    end
  end

  c, res
end

function find_letter(path)
  file  = open(path)
  lines = readlines(file)

  for i = 1:length(lines)
    for j = 1:length(lines)
      if i == j
        continue
      end

      c, res = diff(lines[i], lines[j])

      if c == 1
        return res
      end
    end
  end
end

res2 = find_letter("input2.txt")
println(res2)
