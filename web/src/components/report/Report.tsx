import React, { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { useTournament } from '../../utils/tournament';
import { SetList } from './SetList';
import { SetReporter } from './SetReporter';
import { Set } from '../../types';

export function Report() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [tournament] = useTournament();
  const setIdParam = searchParams.get('set');
  const [selectedSet, setSelectedSet] = useState<Set | null>(null);

  useEffect(() => {
    if (!setIdParam) {
      setSelectedSet(null);
    }
  }, [setIdParam]);

  const handleSelectSet = (set: Set) => {
    setSelectedSet(set);
    setSearchParams({ set: set.id.toString() });
  };

  const handleBack = () => {
    setSearchParams({});
  };

  if (selectedSet) {
    return <SetReporter set={selectedSet} tournament={tournament} onBack={handleBack} />;
  }

  return <SetList tournament={tournament} onSelectSet={handleSelectSet} />;
}
