import { NextApiRequest, NextApiResponse } from 'next';
import DreamRepository from '@/lib/repositories/dream-repository';

export async function POST(request: NextApiRequest, response: NextApiResponse) {
  const body = await request.json();
  const dreamRepository = new DreamRepository();

  try {
    await dreamRepository.create(body.dream);
    // TODO - Update this to modern next
    // return NextResponse.json({ message: 'Dream saved successfully' });
  } catch (error) {
    console.error('Error saving dream:', error);
    // TODO - Update this to modern next
    // return NextResponse.json({ error: 'Failed to save dream' }, { status: 500 });
  }
}

export async function GET(request: NextApiRequest, response: NextApiResponse) {
  const dreamRepository = new DreamRepository();

  try {
    const dreams = await dreamRepository.findAll();
    // TODO - Update this to modern next
    // return NextResponse.json(dreams);
  } catch (error) {
    console.error('Error fetching dreams:', error);
    // TODO - Update this to modern next
    // return NextResponse.json({ error: 'Failed to fetch dreams' }, { status: 500 });
  }
}
