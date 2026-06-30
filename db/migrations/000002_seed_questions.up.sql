-- Starter practice questions aligned with recovery plan task categories:
-- type = programming_basics | dsa | concept (matches plan_tasks.category)

INSERT INTO questions (id, type, category, level, title, content, constraints, starter_code, solution_hint) VALUES
(
    'a1000001-0000-4000-8000-000000000001',
    'programming_basics',
    'loops',
    'level_1',
    'Sum integers from 1 to n',
    'Write a function that returns the sum of all integers from 1 through n (inclusive).

Example: `sumToN(5)` returns `15`.',
    '- n is a non-negative integer
- Use a loop (no closed-form formula required)',
    'function sumToN(n) {
  // your code here
}',
    'Initialize an accumulator to 0 and iterate from 1 to n.'
),
(
    'a1000001-0000-4000-8000-000000000002',
    'programming_basics',
    'functions',
    'level_2',
    'Safe divide with error handling',
    'Implement `divide(a, b)` that returns `{ ok: true, value }` on success or `{ ok: false, error }` when dividing by zero.

Do not throw exceptions.',
    '- a and b are numbers
- Return a structured result object',
    'function divide(a, b) {
  // your code here
}',
    'Check b === 0 before performing division.'
),
(
    'a1000001-0000-4000-8000-000000000003',
    'programming_basics',
    'debugging',
    'level_3',
    'Fix the off-by-one loop',
    'The function below should return the last element of an array but fails on some inputs. Fix the bug.

```javascript
function lastElement(arr) {
  for (let i = 0; i <= arr.length; i++) {
    if (i === arr.length - 1) {
      return arr[i];
    }
  }
  return null;
}
```',
    '- arr is a non-empty array
- Preserve the loop structure; fix the boundary condition',
    NULL,
    'Valid indices for an array of length n are 0 through n - 1.'
),
(
    'a1000001-0000-4000-8000-000000000004',
    'dsa',
    'arrays',
    'level_1',
    'Find maximum in an array',
    'Given an array of integers, return the largest value.

Example: `[3, 1, 4, 1, 5]` returns `5`.',
    '- Array has at least one element
- O(n) time, O(1) extra space',
    'function findMax(nums) {
  // your code here
}',
    'Track a running maximum while scanning the array once.'
),
(
    'a1000001-0000-4000-8000-000000000005',
    'dsa',
    'strings',
    'level_1',
    'Reverse a string',
    'Write a function that reverses a string in place (or returns a new reversed string).

Example: `"hello"` becomes `"olleh"`.',
    '- Input is a non-empty string
- Aim for O(n) time',
    'function reverseString(s) {
  // your code here
}',
    'Two pointers from both ends swapping characters works well.'
),
(
    'a1000001-0000-4000-8000-000000000006',
    'dsa',
    'hash_maps',
    'level_2',
    'Two sum',
    'Given an array of integers `nums` and an integer `target`, return indices of the two numbers that add up to `target`.

You may assume exactly one solution exists.',
    '- 2 <= nums.length <= 10^4
- Return indices in any order
- O(n) average time expected',
    'function twoSum(nums, target) {
  // your code here
}',
    'Store each value''s index in a hash map as you scan.'
),
(
    'a1000001-0000-4000-8000-000000000007',
    'dsa',
    'two_pointer',
    'level_2',
    'Valid palindrome',
    'Given a string, determine if it is a palindrome considering only alphanumeric characters and ignoring case.

Example: `"A man, a plan, a canal: Panama"` is a palindrome.',
    '- Ignore non-alphanumeric characters
- Case-insensitive comparison',
    'function isPalindrome(s) {
  // your code here
}',
    'Use two pointers moving inward after normalizing characters.'
),
(
    'a1000001-0000-4000-8000-000000000008',
    'dsa',
    'recursion',
    'level_3',
    'Fibonacci with memoization',
    'Implement `fib(n)` that returns the nth Fibonacci number (0-indexed: fib(0)=0, fib(1)=1).

Use memoization to avoid exponential recomputation.',
    '- 0 <= n <= 50
- Must handle n=0 and n=1 as base cases',
    'function fib(n) {
  // your code here
}',
    'A top-down recursive approach with a cache map works well.'
),
(
    'a1000001-0000-4000-8000-000000000009',
    'dsa',
    'trees',
    'level_3',
    'Maximum depth of a binary tree',
    'Given the root of a binary tree, return its maximum depth (longest path from root to a leaf, counted in nodes).

Example: tree `[3,9,20,null,null,15,7]` has depth 3.',
    '- Tree has 0 to 10^4 nodes
- Recursive or iterative BFS/DFS both acceptable',
    'function maxDepth(root) {
  // your code here
}',
    'Recursively return 1 + max(left depth, right depth), with null base case 0.'
),
(
    'a1000001-0000-4000-8000-000000000010',
    'concept',
    'complexity',
    'level_1',
    'Big-O of nested loops',
    'What is the time and space complexity of the following function? Explain your reasoning briefly.

```javascript
function countPairs(n) {
  let count = 0;
  for (let i = 0; i < n; i++) {
    for (let j = i; j < n; j++) {
      count++;
    }
  }
  return count;
}
```',
    '- n is a positive integer
- Provide time and space complexity with brief justification',
    NULL,
    'The inner loop runs from i to n-1; sum over i gives O(n²) time. Only O(1) extra space.'
),
(
    'a1000001-0000-4000-8000-000000000011',
    'concept',
    'http',
    'level_2',
    'REST idempotency',
    'Explain which HTTP methods are idempotent and why. Give one example where repeating a non-idempotent request causes a problem.',
    '- Cover GET, PUT, POST, PATCH, DELETE
- One concrete example of duplicate side effects',
    NULL,
    'GET, PUT, DELETE are idempotent; POST typically is not (e.g. double-charging on duplicate POST).'
),
(
    'a1000001-0000-4000-8000-000000000012',
    'concept',
    'sql',
    'level_2',
    'When to add a database index',
    'You have a `users` table with 2M rows. Queries filtering by `email` are slow. Explain when adding an index on `email` helps, and one downside of indexing.',
    '- Assume PostgreSQL
- Mention read vs write trade-offs',
    NULL,
    'B-tree index speeds equality lookups on email; downside is slower inserts/updates and extra storage.'
),
(
    'a1000001-0000-4000-8000-000000000013',
    'programming_basics',
    'functions',
    'level_1',
    'Greet with default name',
    'Write `greet(name)` that returns `"Hello, {name}!"`. If `name` is missing or empty, use `"Guest"`.',
    '- name may be undefined, null, or empty string',
    'function greet(name) {
  // your code here
}',
    'Coalesce falsy/empty name to a default before formatting.'
),
(
    'a1000001-0000-4000-8000-000000000014',
    'dsa',
    'arrays',
    'level_2',
    'Remove duplicates from sorted array',
    'Given a sorted array `nums`, remove duplicates in-place so each element appears at most once. Return the count of unique elements.

Example: `[1,1,2]` → unique count 2, with `nums` beginning `[1,2,...]`.',
    '- Array is sorted ascending
- Modify in-place; O(1) extra space',
    'function removeDuplicates(nums) {
  // your code here
}',
    'Use a write pointer for the next unique position; read pointer scans ahead.'
),
(
    'a1000001-0000-4000-8000-000000000015',
    'concept',
    'complexity',
    'level_3',
    'Hash map vs sorted array lookup',
    'Compare average-case lookup time and memory trade-offs between a hash map and a sorted array for 1M read-heavy key lookups. When would you pick each?',
    '- Assume keys fit in memory
- Mention build/preprocessing cost if relevant',
    NULL,
    'Hash map: O(1) average lookup, higher memory overhead. Sorted array: O(log n) binary search, better cache locality, no hash computation.'
);
