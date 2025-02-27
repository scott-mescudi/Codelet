interface ErrorResp {
  error: string;
  code: number;
}
interface CodeSnippet {
  id: number;
  language: string;
  title: string;
  code: string;
  favorite: boolean;
  private: boolean;
  tags: string[];
  description: string;
}

interface SmallSnippet {
  id: number;
  language: string;
  title: string;
  favorite: boolean;
}

interface CodeSnippetReq {
  language: string;
  title: string;
  code: string;
  favorite?: boolean;
  private?: boolean;
  tags?: string[];
  description?: string;
}

type SmallSnippets = SmallSnippet[];

export async function GetSnippetByID(
  snippetID: number,
  token: string
): Promise<CodeSnippet | undefined> {
  if (!snippetID || snippetID < 0) {
    return undefined;
  }

  if (token === "") {
    return undefined;
  }

  try {
    const resp = await fetch(
      `https://codeletserver-production.up.railway.app/api/v1/user/snippets/${snippetID}`,
      {
        method: "GET",
        headers: {
          Authorization: token,
        },
      }
    );

    if (!resp.ok) {
      const errResp = (await resp.json()) as ErrorResp;
      console.log(errResp);
      return undefined;
    }

    const snippets = (await resp.json()) as CodeSnippet;
    return snippets;
  } catch (err) {
    console.error(err);
    return undefined;
  }
}

export async function getSmallSnippets(
  token: string
): Promise<SmallSnippets | undefined> {
  if (token === "") {
    return undefined;
  }

  try {
    const resp = await fetch(
      "https://codeletserver-production.up.railway.app/api/v1/user/small/snippets",
      {
        method: "GET",
        headers: {
          Authorization: token,
        },
      }
    );

    if (!resp.ok) {
      return undefined;
    }

    const snippets = (await resp.json()) as SmallSnippets;
    if (snippets.length <= 0) {
      return undefined;
    }

    return snippets;
  } catch (err) {
    console.error(err);
    return undefined;
  }
}

export async function DeleteReq(
  token: string,
  id: number
): Promise<boolean | undefined> {
  if (token === "") {
    return undefined;
  }

  try {
    const resp = await fetch(
      `https://codeletserver-production.up.railway.app/api/v1/user/snippets/${id}`,
      {
        method: "DELETE",
        headers: {
          Authorization: token,
        },
      }
    );

    return resp.ok;
  } catch (err) {
    console.error(err);
    return undefined;
  }
}

export async function addSnippetReq(
  token: string,
  language: string,
  title: string,
  tags: string[],
  description: string,
  code: string
): Promise<boolean | undefined> {
  if (token === "" || title === "" || language === "" || code === "") {
    return undefined;
  }

  const body: CodeSnippetReq = {
    language,
    title,
    code,
    favorite: false,
    private: true,
    tags,
    description,
  };

  try {
    const resp = await fetch(
      "https://codeletserver-production.up.railway.app/api/v1/user/snippets",
      {
        method: "POST",
        headers: {
          Authorization: token,
          "Content-Type": "application/json",
        },
        body: JSON.stringify(body),
      }
    );

    if (!resp.ok) {
      console.log(await resp.json());
      return false;
    }

    return true;
  } catch (err) {
    console.error(err);
    return undefined;
  }
}

export async function updateSnippetReq(
  id: number,
  token: string,
  language: string,
  title: string,
  tags: string[],
  description: string,
  code: string
): Promise<boolean | undefined> {
  if (token === "" || title === "" || language === "" || code === "") {
    return undefined;
  }

  const body: CodeSnippetReq = {
    language,
    title,
    code,
    favorite: false,
    private: true,
    tags,
    description,
  };

  try {
    const resp = await fetch(
      `https://codeletserver-production.up.railway.app/api/v1/user/snippets/${id}`,
      {
        method: "PUT",
        headers: {
          Authorization: token,
          "Content-Type": "application/json",
        },
        body: JSON.stringify(body),
      }
    );

    if (!resp.ok) {
      console.log(await resp.json());
      return false;
    }

    return true;
  } catch (err) {
    console.error(err);
    return undefined;
  }
}
