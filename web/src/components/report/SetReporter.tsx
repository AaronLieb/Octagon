import React, { useState, useEffect } from 'react';
import { Set, Character, Game } from '../../types';
import { ListSearch } from './ListSearch';
import API_URL, { getAuthHeaders } from '../../config';
import './SetReporter.css';

interface SetReporterProps {
  set: Set;
  tournament: string;
  onBack: () => void;
}

export function SetReporter({ set, tournament, onBack }: SetReporterProps) {
  const [games, setGames] = useState<Game[]>([]);
  const [p1DefaultChar, setP1DefaultChar] = useState<Character | undefined>();
  const [p2DefaultChar, setP2DefaultChar] = useState<Character | undefined>();
  const [selectingChar, setSelectingChar] = useState<{ player: 1 | 2; gameIndex: number } | null>(null);
  const [error, setError] = useState<string>('');
  const [characters, setCharacters] = useState<Character[]>([]);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    setError('');
    fetch(`${API_URL}/api/characters`, { headers: getAuthHeaders() })
      .then(res => {
        if (!res.ok) throw new Error('Failed to fetch characters');
        return res.json();
      })
      .then(data => {
        const chars = data.characters || [];
        setCharacters(chars);
        
        // Set default characters from cached data if available
        if (set.p1Char) {
          const char = chars.find((c: Character) => c.name === set.p1Char);
          if (char) setP1DefaultChar(char);
        }
        if (set.p2Char) {
          const char = chars.find((c: Character) => c.name === set.p2Char);
          if (char) setP2DefaultChar(char);
        }
      })
      .catch(err => setError(err instanceof Error ? err.message : 'Failed to fetch characters'));
  }, [set]);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (selectingChar) {
        if (e.key === 'Escape') {
          e.preventDefault();
          setSelectingChar(null);
        }
        return;
      }

      if ((e.key === 'a' || e.key === 'j') && !e.shiftKey) {
        e.preventDefault();
        addGame(1);
      } else if ((e.key === 'd' || e.key === 'k') && !e.shiftKey) {
        e.preventDefault();
        addGame(2);
      } else if (e.key === 'A' || e.key === 'q' || e.key === 'J') {
        e.preventDefault();
        setSelectingChar({ player: 1, gameIndex: games.length });
      } else if (e.key === 'D' || e.key === 'e' || e.key === 'K') {
        e.preventDefault();
        setSelectingChar({ player: 2, gameIndex: games.length });
      } else if (e.key === 'Backspace') {
        e.preventDefault();
        if (games.length > 0) {
          setGames(games.slice(0, -1));
          setError('');
        }
      } else if (e.key === 'Enter') {
        e.preventDefault();
        submitSet();
      } else if (e.key === 'Escape') {
        e.preventDefault();
        onBack();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [games, selectingChar, p1DefaultChar, p2DefaultChar]);

  useEffect(() => {
    if (selectingChar) {
      const handleEscape = (e: KeyboardEvent) => {
        if (e.key === 'Escape') {
          e.preventDefault();
          e.stopPropagation();
          setSelectingChar(null);
        }
      };
      window.addEventListener('keydown', handleEscape, true);
      return () => window.removeEventListener('keydown', handleEscape, true);
    }
  }, [selectingChar]);

  const addGame = (winner: 1 | 2) => {
    if (games.length >= 5) return;
    
    const newGame: Game = {
      winner,
      p1Char: p1DefaultChar,
      p2Char: p2DefaultChar,
    };
    
    setGames([...games, newGame]);
    setError('');
  };

  const handleCharSelect = (char: Character) => {
    if (!selectingChar) return;

    const { player, gameIndex } = selectingChar;
    
    if (player === 1) {
      setP1DefaultChar(char);
      setGames(games.map((g, i) => 
        i === gameIndex ? { ...g, p1Char: char } : (!g.p1Char ? { ...g, p1Char: char } : g)
      ));
    } else {
      setP2DefaultChar(char);
      setGames(games.map((g, i) => 
        i === gameIndex ? { ...g, p2Char: char } : (!g.p2Char ? { ...g, p2Char: char } : g)
      ));
    }
    
    setSelectingChar(null);
  };

  const handleGameClick = (gameIndex: number, player: 1 | 2) => {
    addGame(player);
  };

  const validateSet = (): boolean => {
    const p1Wins = games.filter(g => g.winner === 1).length;
    const p2Wins = games.filter(g => g.winner === 2).length;

    const validBo3 = (p1Wins === 2 && p2Wins <= 1) || (p2Wins === 2 && p1Wins <= 1);
    const validBo5 = (p1Wins === 3 && p2Wins <= 2) || (p2Wins === 3 && p1Wins <= 2);

    if (!validBo3 && !validBo5) {
      setError('Invalid score: must be best-of-3 (2-0, 2-1) or best-of-5 (3-0, 3-1, 3-2)');
      return false;
    }

    setError('');
    return true;
  };

  const submitSet = async () => {
    if (!validateSet() || submitting) return;

    setSubmitting(true);
    const gameResults: { winner: number; p1Char: string; p2Char: string }[] = games.map(g => ({
      winner: g.winner,
      p1Char: g.p1Char?.name || '',
      p2Char: g.p2Char?.name || '',
    }));

    try {
      const response = await fetch(`${API_URL}/api/sets/report?tournament=${tournament}`, {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify({ setId: set.id, games: gameResults }),
      });

      if (response.ok) {
        onBack();
      } else {
        setError('Failed to report set');
      }
    } catch (err) {
      setError('Network error');
    } finally {
      setSubmitting(false);
    }
  };

  if (selectingChar) {
    return (
      <div className="set-reporter">
        <h3>Select Character for {selectingChar.player === 1 ? set.player1.name : set.player2.name}</h3>
        <ListSearch
          items={characters}
          getLabel={char => char.name}
          onSelect={handleCharSelect}
          placeholder="Search character..."
        />
      </div>
    );
  }

  const p1Wins = games.filter(g => g.winner === 1).length;
  const p2Wins = games.filter(g => g.winner === 2).length;

  return (
    <div className="set-reporter">
      <div className="header">
        <button onClick={onBack}>← Back</button>
        <h2>{set.round}</h2>
      </div>

      <div className="players-header">
        <div className="player-name">{set.player1.name}</div>
        <div className="score">{p1Wins} - {p2Wins}</div>
        <div className="player-name">{set.player2.name}</div>
      </div>

      <div className="games-grid">
        {Array.from({ length: 5 }).map((_, i) => (
          <div key={i} className="game-row">
            <div
              className={`game-cell ${games[i]?.winner === 1 ? 'win' : games[i] ? 'loss' : ''} ${!games[i] && (p1DefaultChar || p2DefaultChar) ? 'pending' : ''}`}
            >
              <div className="char-icon-btn" onClick={(e) => {
                e.stopPropagation();
                setSelectingChar({ player: 1, gameIndex: i });
              }}>
                {(games[i]?.p1Char || (!games[i] && p1DefaultChar)) ? (
                  <>
                    <div className="char-icon">🎮</div>
                    <div className="char-name">{games[i]?.p1Char?.name || p1DefaultChar?.name}</div>
                  </>
                ) : (
                  <div className="char-placeholder">?</div>
                )}
              </div>
              <div className="game-click-area" onClick={() => handleGameClick(i, 1)}>
                <div className="player-label">{set.player1.name}</div>
              </div>
            </div>
            <div
              className={`game-cell ${games[i]?.winner === 2 ? 'win' : games[i] ? 'loss' : ''} ${!games[i] && (p1DefaultChar || p2DefaultChar) ? 'pending' : ''}`}
            >
              <div className="game-click-area" onClick={() => handleGameClick(i, 2)}>
                <div className="player-label">{set.player2.name}</div>
              </div>
              <div className="char-icon-btn" onClick={(e) => {
                e.stopPropagation();
                setSelectingChar({ player: 2, gameIndex: i });
              }}>
                {(games[i]?.p2Char || (!games[i] && p2DefaultChar)) ? (
                  <>
                    <div className="char-icon">🎮</div>
                    <div className="char-name">{games[i]?.p2Char?.name || p2DefaultChar?.name}</div>
                  </>
                ) : (
                  <div className="char-placeholder">?</div>
                )}
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="controls">
        <div className="hotkeys">
          <span><kbd>A</kbd>/<kbd>J</kbd> P1 Win</span>
          <span><kbd>D</kbd>/<kbd>K</kbd> P2 Win</span>
          <span><kbd>Shift+A</kbd>/<kbd>Q</kbd>/<kbd>Shift+J</kbd> P1 Char</span>
          <span><kbd>Shift+D</kbd>/<kbd>E</kbd>/<kbd>Shift+K</kbd> P2 Char</span>
          <span><kbd>Backspace</kbd> Undo</span>
          <span><kbd>Enter</kbd> Submit</span>
        </div>
        {error && <div className="error">{error}</div>}
        <button onClick={submitSet} className="submit-btn" disabled={submitting}>
          {submitting ? 'Submitting...' : 'Submit Set'}
        </button>
      </div>
    </div>
  );
}
