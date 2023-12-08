class FetchError extends Error {
  info: any;
  status: number;

  constructor(message: string) {
    super(message);
    this.info = null;
    this.status = 0;
  }
}

export const fetcher = async (url: string) => {
  const res = await fetch(url);

  // If the status code is not in the range 200-299,
  // we still try to parse and throw it.
  if (!res.ok) {
    const error = new FetchError("An error occurred while fetching the data.");
    // Attach extra info to the error object.
    error.info = await res.json();
    error.status = res.status;
    throw error;
  }

  return res.json();
};
