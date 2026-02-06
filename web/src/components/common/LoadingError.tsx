import React from 'react';

export function LoadingSpinner() {
  return <div className="loading">Loading...</div>;
}

interface ErrorMessageProps {
  message: string;
}

export function ErrorMessage({ message }: ErrorMessageProps) {
  return <div className="error">Error: {message}</div>;
}
