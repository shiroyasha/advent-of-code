# rubocop:disable all

class Passport
  def initialize(fields)
    @fields = fields
  end

  def valid?
    puts "Validating #{@fields}"

    unless valid_byr?
      puts "invalid byr"
      return false
    end

    unless valid_iyr?
      puts "invalid iyr"
      return false
    end

    unless valid_eyr?
      puts "invalid eyr"
      return false
    end

    unless valid_hgt?
      puts "invalid hgt"
      return false
    end

    unless valid_hcl?
      puts "invalid hcl"
      return false
    end

    unless valid_ecl?
      puts "invalid ecl"
      return false
    end

    unless valid_pid?
      puts "invalid pid"
      return false
    end

    unless valid_cid?
      puts "invalid cid"
      return false
    end

    return true
  end

  def valid_byr?
    val = @fields["byr"]

    val != nil && val.to_i >= 1920 && val.to_i <= 2002
  end

  def valid_iyr?
    val = @fields["iyr"]

    val != nil && val.to_i >= 2010 && val.to_i <= 2020
  end

  def valid_eyr?
    val = @fields["eyr"]

    val != nil && val.to_i >= 2020 && val.to_i <= 2030
  end

  def valid_hgt?
    val = @fields["hgt"]
    return false if val == nil

    metric = val[-2..-1]
    return false unless ["cm", "in"].include?(metric)

    height = val[0..-3].to_i
    if metric == "cm"
      return height >= 150 && height <= 193
    end

    if metric == "in"
      return height >= 59 && height <= 76
    end
  end

  def valid_hcl?
    val = @fields["hcl"]

    return false if val == nil
    return false if val[0] != "#"

    val[1..-1].length == 6 && val[1..-1] =~ /^[0-9a-f]*$/
  end

  def valid_ecl?
    val = @fields["ecl"]
    return false if val == nil

    ["amb", "blu", "brn", "gry", "grn", "hzl", "oth"].include?(val)
  end

  def valid_pid?
    val = @fields["pid"]
    return false if val == nil

    val.length == 9 && val =~ /^[0-9]*$/
  end

  def valid_cid?
    true
  end
end

input = File.read("input1.txt").split("\n\n")

passports = input.map do |raw|
  Passport.new(raw.split("\n").map { |l| l.split(" ") }.flatten.map { |kw| kw.split(":") }.to_h)
end

valid = passports.select { |p| p.valid? }
puts valid.count
