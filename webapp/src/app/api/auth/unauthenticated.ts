export default function handler(req, res) {
  if (req.method === 'POST') {
    req.status(204).end();
  } else {
    res.setHeader('Allow', ['POST']);
    res.status(405).end(`Method ${req.method} Not Allowed`);
  }
}
