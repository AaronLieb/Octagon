import React, { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { SetList } from './SetList';
import { SetReporter } from './SetReporter';
import { Set } from '../../types';

export function Report() {
  const [searchParams, setSearchParams] = useSearchParams();
  
  const tournamentParam = searchParams.get('tournament') || 'octagon';
  const setIdParam = searchParams.get('set');

  const [selectedSet, setSelectedSet] = useState<Set | null>(null);
  const [tournamentInput, setTournamentInput] = useState(tournamentParam);

  useEffect(() => {
    setTournamentInput(tournamentParam);
  }, [tournamentParam]);

  // Clear selected set when setId is removed from URL
  useEffect(() => {
    if (!setIdParam) {
      setSelectedSet(null);
    }
  }, [setIdParam]);

  const handleSelectSet = (set: Set) => {
    setSelectedSet(set);
    setSearchParams({ tournament: tournamentInput, set: set.id.toString() });
  };

  const handleBack = () => {
    setSearchParams({ tournament: tournamentInput });
  };

  const handleTournamentBlur = () => {
    if (tournamentInput !== tournamentParam) {
      setSearchParams({ tournament: tournamentInput });
    }
  };

  const handleTournamentKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleTournamentBlur();
    }
  };

  if (selectedSet) {
    return <SetReporter set={selectedSet} tournament={tournamentInput} onBack={handleBack} />;
  }

  return (
    <div>
      <div style={{ padding: '20px', maxWidth: '800px', margin: '0 auto' }}>
        <input
          type="text"
          value={tournamentInput}
          onChange={e => setTournamentInput(e.target.value)}
          onBlur={handleTournamentBlur}
          onKeyDown={handleTournamentKeyDown}
          placeholder="Tournament slug"
          style={{ padding: '8px', fontSize: '14px', width: '200px', marginBottom: '10px' }}
        />
      </div>
      <SetList tournament={tournamentInput} onSelectSet={handleSelectSet} />
    </div>
  );
}
