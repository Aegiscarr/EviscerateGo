package structs

type RapidSzResponse struct {
	Albums struct {
		TotalCount int
		Items      []struct {
			Data struct {
				URI     string
				Name    string
				Artists struct {
					Items []struct {
						URI     string
						Profile struct {
							Name string
						}
					}
				}
				CoverArt struct {
					Sources []struct {
						URL    string
						Width  int
						Height int
					}
				}
				Date struct {
					Year int
				}
			}
		}
	}

	Artists struct {
		TotalCount int
		Items      []struct {
			Data struct {
				URI     string
				Profile struct {
					Name string
				}
				Visuals struct {
					AvatarImage struct {
						Sources []struct {
							URL    string
							Width  int
							Height int
						}
					}
				}
			}
		}
	}

	// Podcast episodes, unused
	Episodes struct {
		TotalCount int
		Items      []struct {
			Data struct {
				URI      string
				Name     string
				CoverArt struct {
					Sources []struct {
						URL    string
						Width  int
						Height int
					}
				}
				Duration struct {
					TotalMilliseconds int
				}
				ReleaseDate struct {
					ISOString string
				}
				Podcast struct {
					CoverArt struct {
						Sources []struct {
							URL    string
							Width  int
							Height int
						}
					}
				}
				Description   string
				ContentRating struct {
					Label string
				}
			}
		}
	}

	// Genres, unused
	Genres struct {
		TotalCount int
		Items      []struct {
			Data struct {
				URI   string
				Name  string
				Image struct {
					Sources []struct {
						URL    string
						Width  int
						Height int
					}
				}
			}
		}
	}

	// Playlists, unused
	Playlists struct {
		TotalCount int
		Items      []struct {
			Data struct {
				URI         string
				Name        string
				Description string
				Images      struct {
					Items []struct {
						Sources []struct {
							URL    string
							Width  int
							Height int
						}
					}
				}
				Owner struct {
					Name string
				}
			}
		}
	}

	// Podcasts, unused
	Podcasts struct {
		TotalCount int
		Items      []struct {
			Data struct {
				URI      string
				Name     string
				CoverArt struct {
					Sources []struct {
						URL    string
						Width  int
						Height int
					}
				}
				Type      string
				Publisher struct {
					Name string
				}
				MediaType string
			}
		}
	}

	TopResults struct {
		Items []struct {
			Data struct {
				URI         string
				ID          string
				Name        string
				Description string
				Type        string
				Publisher   struct {
					Name string
				}
				MediaType string
				Date      struct {
					Year int
				}
				CoverArt struct {
					Sources []struct {
						URL    string
						Width  int
						Height int
					}
				}
				AlbumOfTrack struct {
					URI      string
					Name     string
					CoverArt struct {
						Sources []struct {
							URL    string
							Width  int
							Height int
						}
					}
					ID          string
					SharingInfo struct {
						ShareURL string
					}
				}
				Artists struct {
					Items []struct {
						URI     string
						Profile struct {
							Name string
						}
					}
				}
				ContentRating struct {
					Label string
				}
				Duration struct {
					TotalMilliseconds int
				}
				Playability struct {
					Playable bool
				}
				Profile struct {
					Name string
				}
				Visuals struct {
					AvatarImage struct {
						Sources []struct {
							URL    string
							Width  int
							Height int
						}
					}
				}
				Images struct {
					Items []struct {
						Sources []struct {
							URL    string
							Width  int
							Height int
						}
					}
				}
				Owner struct {
					Name string
				}
			}
		}

		Featured []struct {
			Data struct {
				URI         string
				Name        string
				Description string
				Images      struct {
					Items []struct {
						Sources []struct {
							URL    string
							Width  int
							Height int
						}
					}
				}
				Owner struct {
					Name string
				}
			}
		}
	}

	Tracks struct {
		TotalCount int
		Items      []struct {
			Data struct {
				URI          string
				ID           string
				Name         string
				AlbumOfTrack struct {
					URI      string
					Name     string
					CoverArt struct {
						Sources []struct {
							URL    string
							Width  int
							Height int
						}
					}
					ID          string
					SharingInfo struct {
						ShareURL string
					}
				}
				Artists struct {
					Items []struct {
						URI     string
						Profile struct {
							Name string
						}
					}
				}
				ContentRating struct {
					Label string
				}
				Duration struct {
					TotalMilliseconds int
				}
				Playability struct {
					Playable bool
				}
			}
		}
	}

	// users, unused
	Users struct {
		TotalCount int
		Items      []struct {
			Data struct {
				URI         string
				ID          string
				DisplayName string
				Username    string
				Image       struct {
					SmallImageURL string
					LargeImageURL string
				}
			}
		}
	}
}
