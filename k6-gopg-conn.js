import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  scenarios: {
    // Define a scenario that ramps up and down repeatedly
    interrupted_load: {
      executor: 'ramping-vus',
      startVUs: 0,
        stages: [
          { duration: '30s', target: 20 },
          { duration: '2m', target: 20 },
          { duration: '2m', target: 20 },
          { duration: '30s', target: 20 },
          { duration: '2m', target: 20 },
          { duration: '2m', target: 20 },
          { duration: '30s', target: 20 },
          { duration: '2m', target: 20 },
          { duration: '2m', target: 20 },
          { duration: '30s', target: 20 },
          { duration: '2m', target: 20 },
          { duration: '2m', target: 20 },
          { duration: '30s', target: 20 },
          { duration: '2m', target: 20 },
          { duration: '2m', target: 20 },
          { duration: '30s', target: 0 },
        ],
      gracefulRampDown: '0s',
      exec: 'userScenario', // The function to execute for this scenario
    },
  },
  thresholds: {
    // Define any thresholds
  },
};

// Main user scenario logic
export function userScenario() {
  const BASE_URL = 'http://localhost:3001';

  const names = [
  'John', 'Jane', 'Doe', 'Michael', 'Sarah', 'Jessica', 'David', 'Emily', 'Ashley',
  'Brian', 'Kimberly', 'Lisa', 'Michelle', 'Anthony', 'Megan', 'Christopher', 'Amanda',
  'Patricia', 'Matthew', 'Melissa', 'Jason', 'Jennifer', 'Heather', 'Nicole', 'Andrew',
  'Elizabeth', 'Adam', 'Kevin', 'Stephanie', 'Ryan', 'Hannah', 'Jeffrey', 'Laura',
  'Amy', 'Rebecca', 'Brittany', 'Danielle', 'Benjamin', 'Katherine', 'Samantha',
  'Timothy', 'Christina', 'Brandon', 'Allison', 'Amber', 'Courtney', 'Jordan',
  'Mary', 'Daniel', 'Kristen', 'Katie', 'Zachary', 'Victoria', 'Erin', 'Nicholas',
  'Kelly', 'Rachel', 'Aaron', 'Chelsea', 'Charles', 'Lindsey', 'Brooke', 'Joshua',
  'Lauren', 'Caitlin', 'Justin', 'Sara', 'Kayla', 'Alexandra', 'Natalie', 'Mark',
  'Tiffany', 'Sean', 'Anna', 'Tyler', 'Kristin', 'Eric', 'Alexis', 'Kyle', 'Kelsey',
  'Jacqueline', 'Shannon', 'Lindsay', 'Holly', 'Austin', 'Molly', 'Taylor', 'Cassandra',
  'Angela', 'Britney', 'Cody', 'Leah', 'Kaylee', 'Kara', 'Catherine', 'James', 'Meghan',
  'Katelyn', 'Gregory', 'Brianne', 'Alyssa', 'Cory', 'Kaitlyn', 'Ethan', 'Casey', 'Olivia',
  'Steven', 'Paige', 'Jordan', 'Whitney', 'Dylan', 'Haley', 'Spencer', 'Ariel',
  'Christopher', 'Chelsea', 'Hillary', 'Erika', 'Megan', 'Travis', 'Jesse', 'Kendra',
  'Brianna', 'Zachary', 'Julie', 'Patrick', 'Brett', 'Madison', 'Vanessa', 'Gabrielle',
  'Sierra', 'Tara', 'Trevor', 'Kaitlin', 'Katherine', 'Derek', 'Bethany', 'Nathan',
  'Marissa', 'Kristina', 'Sean', 'Krystal', 'Miranda', 'Dustin', 'Brandi', 'Blake',
  'Candice', 'Phillip', 'Chelsey', 'Jenna', 'Garrett', 'Sydney', 'Cameron', 'Jared',
  'Cynthia', 'Alex', 'Scott', 'Jasmine', 'Bryan', 'Morgan', 'Evan', 'Cameron',
  'Kaitlyn', 'Crystal', 'Jordan', 'Monica', 'Brian', 'Veronica', 'Samantha', 'Dana',
  'Cory', 'Courtney', 'Kasey', 'Devin', 'Melanie', 'Erika', 'Katie', 'John', 'Sandra',
  'Brittany', 'Branden', 'Tara', 'Kelli', 'Ian', 'Erin', 'Kellie', 'Ashley', 'Chelsie',
  'Alex', 'Cassandra', 'Candace', 'Allison', 'Bradley', 'Jacqueline', 'Shawn', 'Marie',
  'Sabrina', 'Dylan', 'Briana', 'Hilary', 'Catherine', 'Lindsey', 'Krista', 'Chase',
  'Holly', 'Derrick', 'Adriana', 'Nathan', 'Bailey', 'Chelsea', 'Corey', 'Amanda',
  'Shelby', 'Kayla', 'Kimberly', 'Stephen', 'Samantha', 'Meagan', 'Melinda', 'Jordan',
  'Nicole', 'Kayla', 'Brianna', 'Kenneth', 'Alexandra', 'Joshua', 'Kirsten', 'Brooke',
  'Douglas', 'Jenna', 'Michael', 'Elizabeth', 'Kasey', 'Brandi', 'Whitney', 'Patrick',
  'Kristin', 'Kathleen', 'Tori', 'Justin', 'Natalie', 'Katelyn', 'Hillary', 'Kevin',
  'Victoria', 'Kendra', 'Chelsea', 'Heather', 'Phillip', 'Ashlee', 'Kristen', 'Lindsay',
  'Lacey', 'Kristina', 'Brittney', 'Derek', 'Lindsay', 'Aaron', 'Rachel', 'Chelsea',
  'Jennifer', 'Breanna', 'Blake', 'Ryan', 'Ashley', 'Mallory', 'Kylie', 'Jasmine',
  'Kelli', 'Kelsey', 'Emily', 'Kaitlyn', 'Lauren', 'Kara', 'Jordan', 'Alyssa', 'Shannon',
  'Kaylee', 'Alexandria', 'Lindsey', 'Krista', 'Alicia', 'Kayla', 'Katelyn', 'Taylor',
  'Caitlin', 'Dylan', 'Cody', 'Caleb', 'Alexandra', 'Zachary', 'Amber', 'Alyssa',
  'Rebecca', 'Michael', 'Dillon', 'Jennifer', 'Brett', 'Megan', 'Kathryn', 'Haley',
  'Kelsey', 'Jenna', 'Jenna', 'Spencer', 'Katherine', 'Kara', 'Whitney', 'Rachel',
  'Jordan', 'Heather', 'Victoria', 'Ashley', 'Cassandra', 'Hannah', 'Tyler', 'Jessica',
  'Jesse', 'Kayla', 'Chelsea', 'Jennifer', 'Brittany', 'Taylor', 'Chase', 'Kaitlin',
  'Jordan', 'Matthew', 'Erin', 'Haley', 'Hannah', 'Alexis', 'Alyssa', 'Lauren', 'Briana',
  'Casey', 'Chelsea', 'Morgan', 'John', 'Brittney', 'Dylan', 'Courtney', 'Brittany',
  'Samantha', 'Cameron', 'Amanda', 'Tiffany', 'Joshua', 'Brianna', 'Dillon', 'Morgan',
  'Kayla', 'Samantha', 'Amanda', 'Brooke', 'Katie', 'Alyssa', 'Brianna', 'Katie', 'Kaitlin',
  'Elizabeth', 'Kristen', 'Victoria', 'Blake', 'Brittany', 'Whitney'];

  const randomName = names[Math.floor(Math.random() * names.length)];

  const res = http.get(`${BASE_URL}/search/users?name=${randomName}`);
  check(res, { 'status was 200': (r) => r.status == 200 });

  // Check if the response body can be parsed as JSON
  let users;
  try {
    users = res.json();
  } catch (error) {
    console.error(`Failed to parse response as JSON: ${error}`);
    return;
  }

  // Check if the parsed object has a length property
  if (!Array.isArray(users)) {
    console.error('Response is not an array');
    return;
  }

  // Check if the array is empty
  if (users.length === 0) {
    console.log('No users found');
    return;
  }

  const userID = users[0].id;

  const purchasesRes = http.get(`${BASE_URL}/user/${userID}/purchases`);
  check(purchasesRes, { 'status was 200': (r) => r.status == 200 });

  const cancellationsRes = http.get(`${BASE_URL}/user/${userID}/purchases/cancellations`);
  check(cancellationsRes, { 'status was 200': (r) => r.status == 200 });

  sleep(1); // Control iteration pacing here if necessary
}

export function teardown(data) {
  // Optional: Any teardown logic if necessary when the test ends
  console.log("Teardown completed")
}
