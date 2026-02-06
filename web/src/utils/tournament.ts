import React from 'react';

export function getTournament(): string {
  return localStorage.getItem('tournament') || 'octagon';
}

export function setTournament(tournament: string): void {
  localStorage.setItem('tournament', tournament);
  window.dispatchEvent(new Event('tournament-changed'));
}

export function useTournament(): [string, (tournament: string) => void] {
  const [tournament, setTournamentState] = React.useState(getTournament());

  React.useEffect(() => {
    const handleChange = () => {
      setTournamentState(getTournament());
    };
    window.addEventListener('tournament-changed', handleChange);
    return () => window.removeEventListener('tournament-changed', handleChange);
  }, []);

  return [tournament, setTournament];
}

