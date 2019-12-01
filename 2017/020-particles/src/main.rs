extern crate regex;

use std::fs::File;
use std::io::Read;
use regex::Regex;
use std::cmp::Ordering;

#[derive(Eq, Debug, Clone)]
struct Particle {
    index : usize,
    pos: (i64, i64, i64),
    vel: (i64, i64, i64),
    acc: (i64, i64, i64),
}

impl PartialOrd for Particle {
    fn partial_cmp(&self, other: &Particle) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

impl PartialEq for Particle {
    fn eq(&self, other: &Particle) -> bool {
        self.pos == other.pos && self.vel == other.vel && self.acc == other.acc
    }
}

impl Ord for Particle {
    fn cmp(&self, other : &Particle) -> Ordering {
        let acc_dist_1 = taxi_dist(self.acc);
        let acc_dist_2 = taxi_dist(other.acc);

        if acc_dist_1 != acc_dist_2 {
            return acc_dist_1.cmp(&acc_dist_2);
        }

        let vel_dist_1 = taxi_dist(self.vel);
        let vel_dist_2 = taxi_dist(other.vel);

        if vel_dist_1 != vel_dist_2 {
            return vel_dist_1.cmp(&vel_dist_2);
        }

        let dist_1 = taxi_dist(self.pos);
        let dist_2 = taxi_dist(other.pos);

        dist_1.cmp(&dist_2)
    }
}

impl Particle {
    fn advance(&mut self) {
        self.vel = (self.vel.0 + self.acc.0, self.vel.1 + self.acc.1, self.vel.2 + self.acc.2);
        self.pos = (self.pos.0 + self.vel.0, self.pos.1 + self.vel.1, self.pos.2 + self.vel.2);
    }
}

fn parse(filename : &str) -> Vec<Particle> {
    let mut file = File::open(filename).expect("Can't open file");
    let mut content = String::new();

    file.read_to_string(&mut content).expect("Can't read from file");

    let regex = Regex::new(r"p=<(.*),(.*),(.*)>, v=<(.*),(.*),(.*)>, a=<(.*),(.*),(.*)>").unwrap();

    content.lines().enumerate().map(|(index, line)| {
        let captures = regex.captures(line).unwrap();
        let fetch_num = |index| { captures.get(index).unwrap().as_str().trim().parse().unwrap() };

        Particle {
            index : index,
            pos : (fetch_num(1), fetch_num(2), fetch_num(3)),
            vel : (fetch_num(4), fetch_num(5), fetch_num(6)),
            acc : (fetch_num(7), fetch_num(8), fetch_num(9)),
        }
    }).collect()
}

fn taxi_dist(point : (i64, i64, i64)) -> i64 {
    point.0.abs() + point.1.abs() + point.2.abs()
}

fn closest(filename : &str) -> usize {
    let particles = parse(filename);

    particles.iter().min_by(|p1, p2| p1.cmp(&p2)).unwrap().index
}

fn simulate(filename : &str) -> usize {
    let mut particles = parse(filename);
    let mut count = 0;
    let mut non_changed_count = 0;

    loop {
        println!("{}, {}", count, non_changed_count);

        for p in particles.iter_mut() {
            p.advance();
        }

        let mut new_particles = Vec::new();

        for p in particles.iter() {
            let collided = particles.iter().filter(|other| p.index != other.index).any(|other| { other.pos == p.pos });

            if !collided {
                new_particles.push(p.clone());
            }
        }

        particles = new_particles;

        if count == particles.len() {
            non_changed_count += 1;

            if non_changed_count > 500 {
                break;
            }
        } else {
           non_changed_count = 0;
           count = particles.len();
        }
    }

    count
}

fn main() {
    println!("{:?}", simulate("input.txt"));
}

#[test]
fn parse_test() {
    assert_eq!(closest("test_input.txt"), 0);
}
