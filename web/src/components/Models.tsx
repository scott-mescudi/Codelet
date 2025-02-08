interface CodeSnippet {
  language: string;
  title: string;
  code: string;
  favorite: boolean;
  private: boolean;
  tags: string[];
  description: string;
}


interface ErrorResponse {
  error: string;
  code: number;
}

type CodeSnippets = CodeSnippet[];
