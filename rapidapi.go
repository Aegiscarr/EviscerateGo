package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

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

func GetRapidAPICall(parsedQ string, typeQ string) *RapidSzResponse {

	var (
		szresponse  *RapidSzResponse
		parsedQuery = parsedQ
		queryType   = typeQ
		url         string
	)

	url = "https://spotify23.p.rapidapi.com/search/?q=" + parsedQuery + "&type=" + queryType + "&offset=0&limit=1&numberOfTopResults=1"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", *RapidSzToken)
	req.Header.Add("X-RapidAPI-Host", "spotify23.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("Error while sending request: ", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	//fmt.Println(string(body))
	fmt.Println(json.Valid(body))
	fmt.Println(json.Unmarshal(body, &szresponse))

	return szresponse
}
